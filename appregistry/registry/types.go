package registry

type AppEntry struct {
	Name         string `json:"appName,omitempty"`
	Description  string `json:"appDescription,omitempty"`
	Ignite       string `json:"ignite,omitempty"`
	Dependencies struct {
		Docker string `json:"docker,omitempty"`
	} `json:"dependencies,omitempty"`
	CosmosSDK string `json:"cosmosSDK,omitempty"`
	Authors   []struct {
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		Website string `json:"website,omitempty"`
	} `json:"authors,omitempty"`
	RepositoryURL    string `json:"repositoryUrl,omitempty"`
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
