package main

import (
	"github.com/spf13/viper"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/nodes"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/openchain"
	"gitlab.com/fireblocks/web3/utils/evm-cli/cmd"
	decompiler "gitlab.com/fireblocks/web3/utils/evm-cli/decompiler"
)

var (
	client           = nodes.NewEthereumNodeClient("https://eth.llamarpc.com")
	decompilerClient = decompiler.NewDecompiler(openchain.NewClient())
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	cmd.Execute(client, decompilerClient)
}
