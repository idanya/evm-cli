package cmd

import (
	"fmt"
	"os"

	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "evm-cli",
	Short: "A CLI tool to interact with the EVM blockchains via JSON-RPC",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(decompiler *decompiler.Decompiler) {

	rootCmd.PersistentFlags().UintP("chain-id", "c", 1, "Chain ID of the blockchain")
	rootCmd.PersistentFlags().String("rpc-url", "", "node RPC endpoint (overrides the chain ID)")

	viper.BindPFlag("chainId", rootCmd.PersistentFlags().Lookup("chain-id"))
	viper.BindPFlag("rpcUrl", rootCmd.PersistentFlags().Lookup("rpc-url"))

	tx := NewTransactionCommands(decompiler)
	rootCmd.AddCommand(tx.GetRootCommand())

	accountCmd := NewAccountCommands()
	rootCmd.AddCommand(accountCmd.GetRootCommand())

	contractCmd := NewContractCommands(decompiler)
	rootCmd.AddCommand(contractCmd.GetRootCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
