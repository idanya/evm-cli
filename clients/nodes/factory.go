package nodes

var (
	ChainRpc = map[uint]string{
		1: "https://eth.llamarpc.com",
		5: "https://ethereum-goerli.publicnode.com",
	}
)

func NodeClientFactory(chainId uint, rpcUrl string) *EthereumNodeClient {
	if rpcUrl != "" {
		return NewEthereumNodeClient(rpcUrl)
	} else if rpc, ok := ChainRpc[chainId]; ok {
		return NewEthereumNodeClient(rpc)
	}
	panic("rpcUrl is not provided")
}
