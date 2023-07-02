package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/nodes"
	decompiler "gitlab.com/fireblocks/web3/utils/evm-cli/decompiler"
)

var rootCmd = &cobra.Command{
	Use:   "evm-cli",
	Short: "A CLI tool to interact with the EVM blockchains via JSON-RPC",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute(client nodes.NodeClient, decompiler *decompiler.Decompiler) {

	rootCmd.PersistentFlags().Uint("chain-id", 1, "Chain ID of the blockchain")

	tx := NewTransactionCommands(nodes.NodeClientFactory)
	rootTxCmd := tx.GetRootCommand()
	rootTxCmd.AddCommand(tx.GetTransactionDataCommand())
	rootTxCmd.AddCommand(tx.GetTransactionReceiptCommand())
	rootCmd.AddCommand(rootTxCmd)

	accountCmd := NewAccountCommands(nodes.NodeClientFactory)
	rootAccountCmd := accountCmd.GetRootCommand()
	rootAccountCmd.AddCommand(accountCmd.GetAccountNonceCommand())
	rootCmd.AddCommand(rootAccountCmd)

	contractCmd := NewContractCommands(nodes.NodeClientFactory, decompiler)
	rootContractCmd := contractCmd.GetRootCommand()
	rootContractCmd.AddCommand(contractCmd.GetContractOpCodeCommand())
	rootContractCmd.AddCommand(contractCmd.GetContractFunctionListCommand())
	rootContractCmd.AddCommand(contractCmd.GetContractExecCommand())
	rootCmd.AddCommand(rootContractCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
