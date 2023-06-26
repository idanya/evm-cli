package etherscan

import (
	"context"

	"gitlab.com/fireblocks/web3/utils/evm-cli/entities"
)

type EtherscanClient interface {
	GetContractABI(context context.Context, contractAddress string) (entities.ABI, error)
}
