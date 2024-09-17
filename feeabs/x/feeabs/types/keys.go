package types

const (
	// ModuleName defines the module name
	ModuleName = "feeabs"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_feeabs"
)

var (
	ParamsKey = []byte("p_feeabs")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
