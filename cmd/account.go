package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/idanya/evm-cli/clients/nodes"
	"github.com/spf13/cobra"
)

type AccountCommands struct {
	nodeClient nodes.NodeClient
}

func NewAccountCommands(nodeClient nodes.NodeClient) *AccountCommands {
	return &AccountCommands{nodeClient}
}

func (ac *AccountCommands) GetRootCommand() *cobra.Command {
	command := &cobra.Command{
		Use:   "account",
		Short: "Account related commands",
	}
	command.AddCommand(ac.GetAccountNonceCommand())

	return command
}

func (ac *AccountCommands) GetAccountNonceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "nonce <address>",
		Short: "Get account nonce",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			count, err := ac.nodeClient.GetAccountNonce(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(count)
		},
	}
}
