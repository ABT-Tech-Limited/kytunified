package kyt

// Chain ID constants for supported blockchains.
// These values are compatible with Beosin and other KYT providers.

// Full query support chains - these chains support complete KYT functionality
const (
	// ChainIDBTC is the chain ID for Bitcoin
	ChainIDBTC = "0"

	// ChainIDETH is the chain ID for Ethereum Mainnet
	ChainIDETH = "1"

	// ChainIDBSC is the chain ID for BNB Smart Chain (formerly Binance Smart Chain)
	ChainIDBSC = "56"

	// ChainIDTron is the chain ID for Tron
	ChainIDTron = "tron"

	// ChainIDPolygon is the chain ID for Polygon (Matic)
	ChainIDPolygon = "137"

	// ChainIDSolana is the chain ID for Solana
	ChainIDSolana = "solana"

	// ChainIDTON is the chain ID for TON (The Open Network)
	ChainIDTON = "ton"

	// ChainIDArbitrum is the chain ID for Arbitrum One
	ChainIDArbitrum = "42161"

	// ChainIDOptimism is the chain ID for Optimism
	ChainIDOptimism = "10"

	// ChainIDAvalanche is the chain ID for Avalanche C-Chain
	ChainIDAvalanche = "43114"

	// ChainIDFantom is the chain ID for Fantom Opera
	ChainIDFantom = "250"

	// ChainIDCronos is the chain ID for Cronos
	ChainIDCronos = "25"

	// ChainIDzkSync is the chain ID for zkSync Era
	ChainIDzkSync = "324"
)

// Basic query support chains - these chains support limited KYT functionality
const (
	// ChainIDBase is the chain ID for Base
	ChainIDBase = "8453"

	// ChainIDLinea is the chain ID for Linea
	ChainIDLinea = "59144"

	// ChainIDScroll is the chain ID for Scroll
	ChainIDScroll = "534352"

	// ChainIDSui is the chain ID for Sui
	ChainIDSui = "sui"

	// ChainIDSonic is the chain ID for Sonic
	ChainIDSonic = "sonic"
)

// ChainInfo contains information about a blockchain.
type ChainInfo struct {
	// ID is the chain identifier used in API requests.
	ID string

	// Name is the human-readable chain name.
	Name string

	// Symbol is the native token symbol.
	Symbol string

	// FullSupport indicates whether this chain has full KYT support.
	FullSupport bool
}

// SupportedChains returns information about all supported blockchains.
func SupportedChains() []ChainInfo {
	return []ChainInfo{
		// Full support chains
		{ID: ChainIDBTC, Name: "Bitcoin", Symbol: "BTC", FullSupport: true},
		{ID: ChainIDETH, Name: "Ethereum", Symbol: "ETH", FullSupport: true},
		{ID: ChainIDBSC, Name: "BNB Smart Chain", Symbol: "BNB", FullSupport: true},
		{ID: ChainIDTron, Name: "Tron", Symbol: "TRX", FullSupport: true},
		{ID: ChainIDPolygon, Name: "Polygon", Symbol: "MATIC", FullSupport: true},
		{ID: ChainIDSolana, Name: "Solana", Symbol: "SOL", FullSupport: true},
		{ID: ChainIDTON, Name: "TON", Symbol: "TON", FullSupport: true},
		{ID: ChainIDArbitrum, Name: "Arbitrum One", Symbol: "ETH", FullSupport: true},
		{ID: ChainIDOptimism, Name: "Optimism", Symbol: "ETH", FullSupport: true},
		{ID: ChainIDAvalanche, Name: "Avalanche C-Chain", Symbol: "AVAX", FullSupport: true},
		{ID: ChainIDFantom, Name: "Fantom Opera", Symbol: "FTM", FullSupport: true},
		{ID: ChainIDCronos, Name: "Cronos", Symbol: "CRO", FullSupport: true},
		{ID: ChainIDzkSync, Name: "zkSync Era", Symbol: "ETH", FullSupport: true},
		// Basic support chains
		{ID: ChainIDBase, Name: "Base", Symbol: "ETH", FullSupport: false},
		{ID: ChainIDLinea, Name: "Linea", Symbol: "ETH", FullSupport: false},
		{ID: ChainIDScroll, Name: "Scroll", Symbol: "ETH", FullSupport: false},
		{ID: ChainIDSui, Name: "Sui", Symbol: "SUI", FullSupport: false},
		{ID: ChainIDSonic, Name: "Sonic", Symbol: "S", FullSupport: false},
	}
}

// FullSupportChainIDs returns a list of chain IDs with full KYT support.
func FullSupportChainIDs() []string {
	return []string{
		ChainIDBTC,
		ChainIDETH,
		ChainIDBSC,
		ChainIDTron,
		ChainIDPolygon,
		ChainIDSolana,
		ChainIDTON,
		ChainIDArbitrum,
		ChainIDOptimism,
		ChainIDAvalanche,
		ChainIDFantom,
		ChainIDCronos,
		ChainIDzkSync,
	}
}

// IsValidChainID checks if the given chain ID is supported.
func IsValidChainID(chainID string) bool {
	for _, chain := range SupportedChains() {
		if chain.ID == chainID {
			return true
		}
	}
	return false
}

// GetChainInfo returns information about a chain by its ID.
// Returns nil if the chain is not found.
func GetChainInfo(chainID string) *ChainInfo {
	for _, chain := range SupportedChains() {
		if chain.ID == chainID {
			return &chain
		}
	}
	return nil
}
