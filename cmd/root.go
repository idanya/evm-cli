package cmd

import (
	"fmt"
	"os"

	"github.com/idanya/evm-cli/clients/directory"
	"github.com/idanya/evm-cli/clients/nodes"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/idanya/evm-cli/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func Execute(directoryClient directory.DirectoryClient, decompiler *decompiler.Decompiler, decoder *services.Decoder) {

	rootCmd := &cobra.Command{
		Use:   "evm-cli",
		Short: "A CLI tool to interact with the EVM blockchains via JSON-RPC",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rootCmd.PersistentFlags().UintP("chain-id", "c", 1, "Chain ID of the blockchain")
	rootCmd.PersistentFlags().String("rpc-url", "", "node RPC endpoint (overrides the chain ID)")

	viper.BindPFlag("chainId", rootCmd.PersistentFlags().Lookup("chain-id"))
	viper.BindPFlag("rpcUrl", rootCmd.PersistentFlags().Lookup("rpc-url"))

	nodeGenerator := func() nodes.NodeClient { return NodeClientFromViper() }

	contractService := services.NewContractService(nodeGenerator, decompiler, decoder)
	transactionService := services.NewTransactionService(nodeGenerator, directoryClient, decoder)

	tx := NewTransactionCommands(transactionService)
	rootCmd.AddCommand(tx.GetRootCommand())

	accountCmd := NewAccountCommands(nodeGenerator)
	rootCmd.AddCommand(accountCmd.GetRootCommand())

	contractCmd := NewContractCommands(contractService, decompiler, decoder)
	rootCmd.AddCommand(contractCmd.GetRootCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
