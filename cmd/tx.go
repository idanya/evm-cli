package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/idanya/evm-cli/clients/directory/openchain"
	"github.com/idanya/evm-cli/services"
	"github.com/spf13/cobra"
)

type TransactionCommands struct {
}

func NewTransactionCommands() *TransactionCommands {
	return &TransactionCommands{}
}

func (tx *TransactionCommands) GetRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "tx",
		Short: "Transaction related commands",
	}

	command.AddCommand(tx.GetTransactionDataCommand())
	command.AddCommand(tx.GetTransactionReceiptCommand())

	return command
}

func (tx *TransactionCommands) GetTransactionDataCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "info <txHash>",
		Short: "Get transaction data by hash",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			tx, err := NodeClientFromViper().GetTransactionByHash(context.Background(), args[0])
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
			transactionService := services.NewTransactionService(NodeClientFromViper(), openchain.NewClient())
			receipt, err := transactionService.GetTransactionReceipt(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(receipt, "", "  ")
			fmt.Println(string(data))
		},
	}
}
