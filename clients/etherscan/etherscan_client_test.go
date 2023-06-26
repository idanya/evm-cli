package etherscan

import (
	"context"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContractABI(t *testing.T) {
	client := NewGenericEtherscanClient("https://api.etherscan.io/api", "YourApiKeyToken")
	abi, err := client.GetContractABI(context.Background(), "0xdac17f958d2ee523a2206206994597c13d831ec7")

	assert.Nil(t, err)
	assert.NotNil(t, abi)
}
