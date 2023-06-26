package nodes

var (
	ChainRpc = map[uint]string{
		1: "https://eth.llamarpc.com",
		5: "https://ethereum-goerli.publicnode.com",
	}
)

type NodeClientFactoryFunc = func(chainId uint) *EthereumNodeClient

func NodeClientFactory(chainId uint) *EthereumNodeClient {
	return NewEthereumNodeClient(ChainRpc[chainId])
}
