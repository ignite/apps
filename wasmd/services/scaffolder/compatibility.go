package scaffolder

type cosmosVersion string
type wasmdVersion []string

var (
	// this map contains the compatibility versions, for some versions we could have several patched and
	// that could be added here
	compatility = map[cosmosVersion]wasmdVersion{
		cosmosVersion("v0.47.3"): wasmdVersion(
			[]string{
				"v0.44.0",
			},
		),
	}
)
