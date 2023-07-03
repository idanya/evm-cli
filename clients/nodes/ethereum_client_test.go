package nodes

import (
	"context"
	"math/big"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTransaction(t *testing.T) {
	client := NewEthereumNodeClient("https://eth.llamarpc.com")
	tx, err := client.GetTransactionByHash(context.Background(), "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a")
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func TestGetTransactionReceipt(t *testing.T) {
	client := NewEthereumNodeClient("https://eth.llamarpc.com")
	tx, err := client.GetTransactionReceipt(context.Background(), "0x15f8e5ea1079d9a0bb04a4c58ae5fe7654b5b2b4463375ff7ffb490aa0032f3a")
	assert.Nil(t, err)
	assert.NotNil(t, tx)
}

func TestExecuteReadFunction(t *testing.T) {
	client := NewEthereumNodeClient("https://eth.llamarpc.com")
	balances, err := client.ExecuteReadFunction(context.Background(), "0xdAC17F958D2ee523a2206206994597C13D831ec7", []string{"address"}, []string{"uint256"}, "balances", "0xdac17f958d2ee523a2206206994597c13d831ec7")
	assert.Nil(t, err)
	assert.NotNil(t, balances)
	assert.Len(t, balances, 1)

	value, _ := new(big.Int).SetString("1384361901572", 10)
	assert.Equal(t, value, balances[0])
}

func TestExecuteReadFunctionNoParams(t *testing.T) {
	client := NewEthereumNodeClient("https://eth.llamarpc.com")
	name, err := client.ExecuteReadFunction(context.Background(), "0xdAC17F958D2ee523a2206206994597C13D831ec7", []string{}, []string{"string"}, "name")
	assert.Nil(t, err)
	assert.Len(t, name, 1)
	assert.Equal(t, "Tether USD", name[0])

	// hex.EncodeToString(name)

	// assert.Equal(t, "Tether USD", string(nameBytes))
}
