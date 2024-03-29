package schema

import "fmt"

const (
	// oracle chainType
	OracleEthChainType     = "ethereum"
	OracleMoonChainType    = "moon"
	OracleCfxChainType     = "conflux"
	OracleBscChainType     = "bsc"
	OraclePlatonChainType  = "platon"
	OracleArweaveChainType = "arweave"
	OracleEverpayChainType = "everpay"

	// ever chainType
	ChainTypeArweave    = "arweave"
	ChainTypeCrossArEth = "arweave,ethereum" // cross arweave, Only used AR token
	ChainTypeMoonbeam   = "moonbeam"
	ChainTypeMoonbase   = "moonbase"
	ChainTypeEth        = "ethereum"
	ChainTypeCfx        = "conflux"
	ChainTypeBsc        = "bsc"
	ChainTypePlaton     = "platon"
	ChainTypeEverpay    = "everpay"
)

func GetEverToNativeChainType(everChainType string) (string, error) {
	switch everChainType { // todo more chain.
	case ChainTypeArweave, ChainTypeCrossArEth:
		return OracleArweaveChainType, nil
	case ChainTypeEth:
		return OracleEthChainType, nil
	case ChainTypeMoonbeam, ChainTypeMoonbase:
		return OracleMoonChainType, nil
	case ChainTypeCfx:
		return OracleCfxChainType, nil
	case ChainTypeBsc:
		return OracleBscChainType, nil
	case ChainTypePlaton:
		return OraclePlatonChainType, nil
	case ChainTypeEverpay:
		return OracleEverpayChainType, nil
	default:
		return "", fmt.Errorf("not found this everChainType:%s", everChainType)
	}
}

type evmChainMeta struct {
	Symbol string
}

var EvmChainMetaMap = map[string]evmChainMeta{ // key: oracle chainType, val: chainInfo
	OracleEthChainType:    {Symbol: "ETH"},  // "ethereum"
	OracleMoonChainType:   {Symbol: "GLMR"}, // "moon"
	OracleCfxChainType:    {Symbol: "CFX"},  // "conflux"
	OracleBscChainType:    {Symbol: "BNB"},  // "bsc"
	OraclePlatonChainType: {Symbol: "LAT"},  // "platon"
}

func IsEvmChain(chainType string) bool { // todo more chain.
	_, ok := EvmChainMetaMap[chainType]
	return ok
}
