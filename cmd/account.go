package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/idanya/evm-cli/entities"
	"github.com/spf13/cobra"
)

type AccountCommands struct {
	nodeClientGenerator entities.NodeClientGenerator
}

func NewAccountCommands(nodeClientGenerator entities.NodeClientGenerator) *AccountCommands {
	return &AccountCommands{nodeClientGenerator}
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
			count, err := ac.nodeClientGenerator().GetAccountNonce(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(count)
		},
	}
}
