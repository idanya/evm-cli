package main

import (
	"github.com/idanya/evm-cli/clients/directory/openchain"
	"github.com/idanya/evm-cli/cmd"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/idanya/evm-cli/services"
	"github.com/spf13/viper"
)

var (
	openChainClient  = openchain.NewClient()
	decompilerClient = decompiler.NewDecompiler(openChainClient)
	decoder          = services.NewDecoder(openChainClient)
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	cmd.Execute(openChainClient, decompilerClient, decoder)
}
