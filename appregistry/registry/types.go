package registry

import (
	"github.com/Masterminds/semver"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/iancoleman/strcase"
	"github.com/ignite/cli/v28/ignite/pkg/errors"
	"net/http"
	"net/mail"
	"net/url"
	"strings"
)

type (
	URL     string
	URLs    []URL
	Version string
	Field   string
	Address string
	Email   string

	Author struct {
		Name    Field `json:"name,omitempty"`
		Email   Email `json:"email,omitempty"`
		Website URL   `json:"website,omitempty"`
	}

	License struct {
		Name Field `json:"name,omitempty"`
		URL  URL   `json:"url,omitempty"`
	}

	SocialMedia struct {
		X        URL    `json:"x,omitempty"`
		Telegram string `json:"telegram,omitempty"`
		Discord  string `json:"discord,omitempty"`
		Reddit   string `json:"reddit,omitempty"`
		Website  URL    `json:"website,omitempty"`
	}

	CryptoAddresses struct {
		Cosmos                Address           `json:"cosmos,omitempty"`
		OtherSupportedCryptos map[string]string `json:"otherSupportedCryptos,omitempty"`
	}

	Donations struct {
		CryptoAddresses   CryptoAddresses `json:"cryptoAddresses,omitempty"`
		FiatDonationLinks URLs            `json:"fiatDonationLinks,omitempty"`
	}

	App struct {
		Name               string            `json:"appName,omitempty"`
		Slug               string            `json:"slug,omitempty"`
		Description        string            `json:"appDescription,omitempty"`
		Ignite             Version           `json:"ignite,omitempty"`
		Dependencies       map[string]string `json:"dependencies,omitempty"`
		CosmosSDK          Version           `json:"cosmosSDK,omitempty"`
		Authors            []Author          `json:"authors,omitempty"`
		RepositoryURL      URL               `json:"repositoryUrl,omitempty"`
		DocumentationURL   URL               `json:"documentationUrl,omitempty"`
		License            License           `json:"license,omitempty"`
		Keywords           []string          `json:"keywords,omitempty"`
		SupportedPlatforms []string          `json:"supportedPlatforms,omitempty"`
		SocialMedia        SocialMedia       `json:"socialMedia,omitempty"`
		Donations          Donations         `json:"donations,omitempty"`
		Icon               URL               `json:"icon,omitempty"`
		Cover              URL               `json:"cover,omitempty"`
	}
)

type (
	FieldCase int

	ValidateOption struct {
		required  bool
		fieldCase FieldCase
		minLength int
	}

	ValidateOptions func(o *ValidateOption)
)

func ValidateFieldCase(fieldCase FieldCase) ValidateOptions {
	return func(f *ValidateOption) {
		f.fieldCase = fieldCase
	}
}

func ValidateRequired() ValidateOptions {
	return func(f *ValidateOption) {
		f.required = true
	}
}

func ValidationLength(minLength int) ValidateOptions {
	return func(f *ValidateOption) {
		f.minLength = minLength
	}
}

const (
	NoCaseSensitive FieldCase = iota
	IsLowerCamel
	IsUpperCamel
	IsLowerCase
	IsUpperCase
	IsKebabCase
	IsSnakeCase
)

func (f Field) ValidateField(opts ...ValidateOptions) error {
	o := ValidateOption{fieldCase: NoCaseSensitive}
	for _, opt := range opts {
		opt(&o)
	}

	field := string(f)
	if field == "" {
		if o.required {
			return errors.Errorf("field %s is required", field)
		}
		return nil
	}

	if o.minLength > 0 && len(field) < o.minLength {
		return errors.Errorf("field %s must be at least %d characters long", field, o.minLength)
	}

	switch o.fieldCase {
	case NoCaseSensitive:
		break
	case IsLowerCamel:
		if field != strcase.ToLowerCamel(field) {
			return errors.Errorf("field %s must be in lower camel case", field)
		}
	case IsUpperCamel:
		if field != strcase.ToCamel(field) {
			return errors.Errorf("field %s must be in upper camel case", field)
		}
	case IsLowerCase:
		if field != strings.ToLower(field) {
			return errors.Errorf("field %s must be in lower case", field)
		}
	case IsUpperCase:
		if field != strings.ToUpper(field) {
			return errors.Errorf("field %s must be in upper case", field)
		}
	case IsKebabCase:
		if field != strcase.ToKebab(field) {
			return errors.Errorf("field %s must be in kebab case", field)
		}
	case IsSnakeCase:
		if field != strcase.ToSnake(field) {
			return errors.Errorf("field %s must be in snake case", field)
		}
	}

	return nil
}

func (v Version) Validate() error {
	_, err := semver.NewConstraint(string(v))
	if err != nil {
		return errors.Wrapf(err, "invalid version %s", v)
	}
	return nil
}

func (e Email) Validate() error {
	if e == "" {
		return errors.Errorf("email must be defined")
	}
	addr, err := mail.ParseAddress(string(e))
	if err != nil {
		return errors.Errorf("invalid email format: %s", err)
	}
	if addr.Address != string(e) {
		return errors.Errorf("invalid email format")
	}
	return nil
}

func (v Version) Verify(version string) error {
	if version == "" {
		return errors.Errorf("version must be defined")
	}
	c, err := semver.NewConstraint(string(v))
	if err != nil {
		return errors.Wrapf(err, "invalid version %s", v)
	}

	semVersion, err := semver.NewVersion(version)
	if err != nil {
		return errors.Wrapf(err, "invalid version %s", v)
	}
	if !c.Check(semVersion) {
		return errors.Errorf("version %s does not satisfy %s", version, v)
	}
	return nil
}

func (a Address) Validate() error {
	if a == "" {
		return errors.Errorf("address must be defined")
	}
	_, _, err := bech32.DecodeAndConvert(string(a))
	if err != nil {
		return errors.Errorf("invalid address %s", a)
	}
	return nil
}

func (us URLs) Validate() error {
	for _, u := range us {
		if err := u.Validate(); err != nil {
			return errors.Wrapf(err, "invalid url %s", u)
		}
	}
	return nil
}

func (u URL) Validate() error {
	if u == "" {
		return errors.Errorf("%s must be defined", u)
	}
	_, err := url.Parse(string(u))
	if err != nil {
		return errors.Errorf("invalid %s format: %s", u, err)
	}
	resp, err := http.Head(string(u))
	if err != nil {
		return errors.Errorf("unable to reach %s: %s", u, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.Errorf("%s returned status code %d", u, resp.StatusCode)
	}
	return nil
}
