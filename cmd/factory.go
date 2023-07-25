package cmd

import (
	"github.com/idanya/evm-cli/clients/nodes"
	"github.com/spf13/viper"
)


func NodeClientFromViper() *nodes.EthereumNodeClient {
	chainId := viper.GetUint("chainId")
	rpcUrl := viper.GetString("rpcUrl")
	return nodes.NodeClientFactory(chainId, rpcUrl)
}
