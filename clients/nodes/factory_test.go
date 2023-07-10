package nodes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClientFactory(t *testing.T) {
	client := NodeClientFactory(1, "")
	assert.Equal(t, ChainRpc[1], client.rpc)

	client = NodeClientFactory(1, "https://example.com")
	assert.Equal(t, "https://example.com", client.rpc)

	client = NodeClientFactory(5, "")
	assert.Equal(t, ChainRpc[5], client.rpc)
}

func TestNodeClientFactory_Panic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Code should have panic")
		}
	}()

	NodeClientFactory(999, "")
}
