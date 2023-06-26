package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/nodes"
)

type TransactionCommands struct {
	clientFactory nodes.NodeClientFactoryFunc
}

func NewTransactionCommands(clientFactory nodes.NodeClientFactoryFunc) *TransactionCommands {
	return &TransactionCommands{clientFactory}
}

func (tx *TransactionCommands) GetRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "tx",
		Short: "Transaction related commands",
	}
}

func (tx *TransactionCommands) GetTransactionDataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info <txHash>",
		Short: "Get transaction data by hash",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}

			tx, err := tx.clientFactory(chainId).GetTransactionByHash(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(tx, "", "  ")
			fmt.Println(string(data))
		},
	}
}

func (tx *TransactionCommands) GetTransactionReceiptCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "receipt <txHash>",
		Short: "Get transaction receipt by hash",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}

			receipt, err := tx.clientFactory(chainId).GetTransactionReceipt(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(receipt, "", "  ")
			fmt.Println(string(data))
		},
	}
}
