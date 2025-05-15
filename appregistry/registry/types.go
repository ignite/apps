package registry

import (
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"net/url"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ignite/cli/v28/ignite/pkg/errors"

	"github.com/ignite/apps/appregistry/strcase"
)

type (
	Apps         []App
	URL          string
	URLs         []URL
	Version      string
	Field        string
	Fields       []Field
	Address      string
	Email        string
	Authors      []Author
	Dependencies map[string]Version

	// Author represents an author with name, email and website information
	Author struct {
		Name    Field `json:"name,omitempty"`
		Email   Email `json:"email,omitempty"`
		Website URL   `json:"website,omitempty"`
	}

	// License contains name and URL of a license
	License struct {
		Name Field `json:"name,omitempty"`
		URL  URL   `json:"url,omitempty"`
	}

	// SocialMedia contains social media links and profiles
	SocialMedia struct {
		X        string `json:"x,omitempty"`
		Telegram string `json:"telegram,omitempty"`
		Discord  string `json:"discord,omitempty"`
		Reddit   string `json:"reddit,omitempty"`
		Website  URL    `json:"website,omitempty"`
	}

	// CryptoAddresses contains Cosmos address and other crypto addresses
	CryptoAddresses struct {
		Cosmos                Address           `json:"cosmos,omitempty"`
		OtherSupportedCryptos map[string]string `json:"otherSupportedCryptos,omitempty"`
	}

	// Donations contains crypto addresses and fiat donation links
	Donations struct {
		CryptoAddresses   CryptoAddresses `json:"cryptoAddresses,omitempty"`
		FiatDonationLinks URLs            `json:"fiatDonationLinks,omitempty"`
	}

	// App represents an Ignite application with its metadata
	App struct {
		Name               Field        `json:"appName,omitempty"`
		Slug               Field        `json:"slug,omitempty"`
		Description        Field        `json:"appDescription,omitempty"`
		Ignite             Version      `json:"ignite,omitempty"`
		Dependencies       Dependencies `json:"dependencies,omitempty"`
		CosmosSDK          Version      `json:"cosmosSDK,omitempty"`
		Authors            Authors      `json:"authors,omitempty"`
		RepositoryURL      URL          `json:"repositoryUrl,omitempty"`
		DocumentationURL   URL          `json:"documentationUrl,omitempty"`
		License            License      `json:"license,omitempty"`
		Keywords           []string     `json:"keywords,omitempty"`
		SupportedPlatforms []string     `json:"supportedPlatforms,omitempty"`
		SocialMedia        SocialMedia  `json:"socialMedia,omitempty"`
		Donations          Donations    `json:"donations,omitempty"`
		Icon               URL          `json:"icon,omitempty"`
		Cover              URL          `json:"cover,omitempty"`
	}
)

type (
	// FieldCase represents different cases for field validation
	FieldCase int

	// ValidateOption contains validation options for fields
	ValidateOption struct {
		required  bool
		fieldCase FieldCase
		minLength int
	}

	// ValidateOptions is a function type that configures validation options
	ValidateOptions func(o *ValidateOption)
)

const (
	// CaseNoSensitive indicates no case sensitivity check
	CaseNoSensitive FieldCase = iota
	// CaseLowerCamel indicates lowercase camel case check
	CaseLowerCamel
	// CaseUpperCamel indicates uppercase camel case check
	CaseUpperCamel
	// CaseLower indicates lowercase check
	CaseLower
	// CaseUpper indicates uppercase check
	CaseUpper
	// CaseKebab indicates kebab case check
	CaseKebab
	// CaseSnake indicates snake case check
	CaseSnake
)

const appsRepoURL = "github.com/ignite/apps/"

// ValidateFieldCase returns a ValidateOptions that sets the field case validation
func ValidateFieldCase(fieldCase FieldCase) ValidateOptions {
	return func(f *ValidateOption) {
		f.fieldCase = fieldCase
	}
}

// ValidateRequired returns a ValidateOptions that marks a field as required
func ValidateRequired() ValidateOptions {
	return func(f *ValidateOption) {
		f.required = true
	}
}

// ValidationLength returns a ValidateOptions that sets minimum length requirement
func ValidationLength(minLength int) ValidateOptions {
	return func(f *ValidateOption) {
		f.minLength = minLength
	}
}

// Validate validates a Field according to the provided options
func (f Field) Validate(opts ...ValidateOptions) error {
	o := ValidateOption{fieldCase: CaseNoSensitive}
	for _, opt := range opts {
		opt(&o)
	}

	field := f.String()
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
	case CaseNoSensitive:
		break
	case CaseLowerCamel:
		if field != strcase.ToLowerCamel(field) {
			return errors.Errorf("field %s must be in lower camel case", field)
		}
	case CaseUpperCamel:
		if field != strcase.ToUpperCamel(field) {
			return errors.Errorf("field %s must be in upper camel case", field)
		}
	case CaseLower:
		if field != strcase.ToLower(field) {
			return errors.Errorf("field %s must be in lower case", field)
		}
	case CaseUpper:
		if field != strcase.ToUpper(field) {
			return errors.Errorf("field %s must be in upper case", field)
		}
	case CaseKebab:
		if field != strcase.ToKebab(field) {
			return errors.Errorf("field %s must be in kebab case", field)
		}
	case CaseSnake:
		if field != strcase.ToSnake(field) {
			return errors.Errorf("field %s must be in snake case", field)
		}
	}

	return nil
}

// Validate validates a slice of Fields according to the provided options
func (f Fields) Validate(opts ...ValidateOptions) error {
	for _, field := range f {
		if err := field.Validate(opts...); err != nil {
			return errors.Wrapf(err, "invalid field %s", field)
		}
	}
	return nil
}

// Validate validates Author fields
func (a Author) Validate() error {
	if a.Name != "" {
		if err := a.Name.Validate(ValidateRequired()); err != nil {
			return errors.Wrapf(err, "invalid author name %s", a.Name)
		}
	}
	if a.Email != "" {
		if err := a.Email.Validate(); err != nil {
			return errors.Wrapf(err, "invalid author email %s", a.Email)
		}
	}
	if a.Website != "" {
		if err := a.Website.Validate(); err != nil {
			return errors.Wrapf(err, "invalid author website %s", a.Website)
		}
	}
	return nil
}

// Validate validates a slice of Authors
func (a Authors) Validate() error {
	for _, author := range a {
		if err := author.Validate(); err != nil {
			return errors.Wrapf(err, "invalid author %s", author.Name)
		}
	}
	return nil
}

// Validate validates Dependencies
func (d Dependencies) Validate() error {
	for name, dep := range d {
		if err := dep.Validate(); err != nil {
			return errors.Wrapf(err, "invalid dependency %s", name)
		}
	}
	return nil
}

// Validate validates Version format
func (v Version) Validate() error {
	_, err := semver.NewConstraint(string(v))
	if err != nil {
		return errors.Wrapf(err, "invalid version %s", v)
	}
	return nil
}

// Validate validates Email format
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

// Verify checks if the provided version satisfies the version constraint
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

// String returns the string representation of the Version
func (v Version) String() string {
	return string(v)
}

// Validate validates Address format
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

// Validate validates License fields
func (l License) Validate() error {
	if err := l.Name.Validate(ValidateRequired()); err != nil {
		return errors.Wrapf(err, "invalid license name %s", l.Name)
	}
	if err := l.URL.Validate(); err != nil {
		return errors.Wrapf(err, "invalid license URL %s", l.URL)
	}
	return nil
}

// Validate validates URLs
func (us URLs) Validate() error {
	for _, u := range us {
		if err := u.Validate(); err != nil {
			return errors.Wrapf(err, "invalid url %s", u)
		}
	}
	return nil
}

// FindBySlug finds an App by its slug
func (apps Apps) FindBySlug(slug string) (App, error) {
	var appEntry App
	for _, app := range apps {
		if strings.EqualFold(string(app.Slug), slug) {
			appEntry = app
		}
	}

	if appEntry.Name == "" && appEntry.Slug == "" {
		return appEntry, errors.Errorf("app slug %s not found", slug)
	}
	return appEntry, nil
}

// FindByName finds an App by its name
func (apps Apps) FindByName(name string) (App, error) {
	var appEntry App
	for _, app := range apps {
		if strings.EqualFold(string(app.Name), name) {
			appEntry = app
		}
	}

	if appEntry.Name == "" && appEntry.Slug == "" {
		return appEntry, errors.Errorf("app name %s not found", name)
	}
	return appEntry, nil
}

// Validate validates URL format and accessibility
func (u URL) Validate() error {
	if u == "" {
		return errors.Errorf("%s must be defined", u)
	}
	if strings.Contains(string(u), appsRepoURL) &&
		!strings.Contains(string(u), "/tree/") &&
		!strings.Contains(string(u), "/blob/") {
		u = URL(strings.ReplaceAll(string(u), appsRepoURL, appsRepoURL+"tree/main/"))
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

// Validate validates App fields
func (a App) Validate() error {
	if err := a.Name.Validate(ValidateRequired(), ValidateFieldCase(CaseUpperCamel)); err != nil {
		return errors.Wrapf(err, "invalid app name %s", a.Name)
	}

	if err := a.Slug.Validate(ValidateRequired(), ValidateFieldCase(CaseKebab)); err != nil {
		return errors.Wrapf(err, "invalid app slug %s", a.Slug)
	}

	if err := a.Description.Validate(ValidateRequired(), ValidationLength(10)); err != nil {
		return errors.Wrapf(err, "invalid app description %s", a.Description)
	}

	if err := a.Ignite.Validate(); err != nil {
		return errors.Wrapf(err, "invalid ignite version %s", a.Ignite)
	}

	if err := a.Dependencies.Validate(); err != nil {
		return errors.Wrapf(err, "invalid dependencies %s", a.Dependencies)
	}

	if err := a.CosmosSDK.Validate(); err != nil {
		return errors.Wrapf(err, "invalid cosmos sdk version %s", a.CosmosSDK)
	}

	if err := a.Authors.Validate(); err != nil {
		return errors.Wrapf(err, "invalid authors %s", a.Authors)
	}

	if err := a.RepositoryURL.Validate(); err != nil {
		return errors.Wrapf(err, "invalid repository url %s", a.RepositoryURL)
	}

	if err := a.DocumentationURL.Validate(); err != nil {
		return errors.Wrapf(err, "invalid documentation url %s", a.DocumentationURL)
	}

	if err := a.License.Validate(); err != nil {
		return errors.Wrapf(err, "invalid license %s", a.License)
	}

	if len(a.Keywords) == 0 {
		return errors.Errorf("unless one keyword must be defined")
	}

	if len(a.SupportedPlatforms) == 0 {
		return errors.Errorf("unless one supportedPlatforms must be defined")
	}

	if a.SocialMedia.Website != "" {
		if err := a.SocialMedia.Website.Validate(); err != nil {
			return errors.Wrapf(err, "invalid social media website %s", a.SocialMedia)
		}
	}

	if a.Donations.CryptoAddresses.Cosmos != "" {
		if err := a.Donations.CryptoAddresses.Cosmos.Validate(); err != nil {
			return errors.Wrapf(err, "invalid cosmos crypto address %s", a.Donations.CryptoAddresses.Cosmos)
		}
	}

	if len(a.Donations.FiatDonationLinks) > 0 {
		if err := a.Donations.FiatDonationLinks.Validate(); err != nil {
			return errors.Wrap(err, "invalid fiat donation link")
		}
	}

	if a.Icon != "" {
		if err := a.Icon.Validate(); err != nil {
			return errors.Wrapf(err, "invalid icon url %s", a.Icon)
		}
	}

	if a.Cover != "" {
		if err := a.Cover.Validate(); err != nil {
			return errors.Wrapf(err, "invalid cover url %s", a.Cover)
		}
	}

	return nil
}

// String returns the string representation of the Field
func (f Field) String() string {
	return string(f)
}

// String returns the string representation of the URL
func (u URL) String() string {
	return string(u)
}

// AppFromFile reads and parses an App from a JSON file
func AppFromFile(r io.Reader) (*App, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read file content")
	}

	var entry *App
	if err := json.Unmarshal(body, &entry); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal file")
	}
	return entry, nil
}
