package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gitlab.com/fireblocks/web3/utils/evm-cli/clients/nodes"
)

type AccountCommands struct {
	clientFactory nodes.NodeClientFactoryFunc
}

func NewAccountCommands(clientFactory nodes.NodeClientFactoryFunc) *AccountCommands {
	return &AccountCommands{clientFactory}
}

func (ac *AccountCommands) GetRootCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "account",
		Short: "Account related commands",
	}
}

func (ac *AccountCommands) GetAccountNonceCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "nonce <address>",
		Short: "Get account nonce",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			chainId, err := cmd.Flags().GetUint("chain-id")
			if err != nil {
				log.Fatal(err)
			}
			count, err := ac.clientFactory(chainId).GetAccountNonce(context.Background(), args[0])
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(count)
		},
	}
}
