package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/idanya/evm-cli/clients/directory/openchain"
	decompiler "github.com/idanya/evm-cli/decompiler"
	"github.com/idanya/evm-cli/services"
	"github.com/spf13/cobra"
)

type TransactionCommands struct {
	decompiler *decompiler.Decompiler
}

func NewTransactionCommands(decompiler *decompiler.Decompiler) *TransactionCommands {
	return &TransactionCommands{decompiler}
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

			tx, err := tx.getTransactionService().GetTransactionByHash(context.Background(), args[0])
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
			receipt, err := tx.getTransactionService().GetTransactionReceipt(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			data, _ := json.MarshalIndent(receipt, "", "  ")
			fmt.Println(string(data))
		},
	}
}

func (tx *TransactionCommands) getTransactionService() *services.TransactionService {
	openchainClient := openchain.NewClient()
	contractService := services.NewContractService(NodeClientFromViper(), tx.decompiler, openchainClient)
	return services.NewTransactionService(NodeClientFromViper(), openchainClient, contractService)
}
