package main

import (
	"github.com/idanya/evm-cli/clients/openchain"
	"github.com/idanya/evm-cli/cmd"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/spf13/viper"
)

var (
	decompilerClient = decompiler.NewDecompiler(openchain.NewClient())
)

func init() {
	viper.AutomaticEnv()
}

func main() {
	cmd.Execute(decompilerClient)
}
