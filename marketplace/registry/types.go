package registry

type AppEntry struct {
	Name         string `json:"appName,omitempty"`
	Description  string `json:"appDescription,omitempty"`
	IgniteVerion string `json:"ignite,omitempty"`
	Dependencies struct {
		Docker string `json:"docker,omitempty"`
	} `json:"dependencies,omitempty"`
	CosmosSDKVersion string   `json:"cosmosSDK,omitempty"`
	Features         []string `json:"features,omitempty"`
	Wasm             bool     `json:"wasm,omitempty"`
	Authors          []struct {
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		Website string `json:"website,omitempty"`
	} `json:"authors,omitempty"`
	Repository struct {
		URL string `json:"url,omitempty"`
	} `json:"repository,omitempty"`
	DocumentationURL string `json:"documentationUrl,omitempty"`
	License          struct {
		Name string `json:"name,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"license,omitempty"`
	Keywords           []string `json:"keywords,omitempty"`
	SupportedPlatforms []string `json:"supportedPlatforms,omitempty"`
	SocialMedia        struct {
		X        string `json:"x,omitempty"`
		Telegram string `json:"telegram,omitempty"`
		Discord  string `json:"discord,omitempty"`
		Reddit   string `json:"reddit,omitempty"`
		Website  string `json:"website,omitempty"`
	} `json:"socialMedia,omitempty"`
	Donations struct {
		CryptoAddresses struct {
			Cosmos                string `json:"cosmos,omitempty"`
			OtherSupportedCryptos string `json:"otherSupportedCryptos,omitempty"`
		} `json:"cryptoAddresses,omitempty"`
		FiatDonationLinks string `json:"fiatDonationLinks,omitempty"`
	} `json:"donations,omitempty"`
}
