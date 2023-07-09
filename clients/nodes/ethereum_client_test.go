package nodes

import (
	"context"

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
