package openchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupFunction(t *testing.T) {
	client := NewClient()
	response, err := client.LookupFunction("0x8b477adb")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "transferFromByLegacy(address,address,address,uint256)", response.Name)
}
