package openchain

import (
	"encoding/json"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"
)

func TestLookupFunction(t *testing.T) {
	defer gock.Off() // Flush pending mocks after test execution
	funcHash := "0x8b477adb"

	expectedResponse := `{"ok":true,"result":{"event":{},"function":{"0x8b477adb":[{"name":"transferFromByLegacy(address,address,address,uint256)","filtered":false}]}}}`
	var expectedResponseJson map[string]interface{}
	err := json.Unmarshal([]byte(expectedResponse), &expectedResponseJson)
	assert.Nil(t, err)

	gock.New("https://api.openchain.xyz").
		Get("/signature-database/v1/lookup").
		MatchParams(map[string]string{"function": funcHash, "filter": "true"}).
		Reply(200).
		JSON(expectedResponseJson)

	client := NewClient()

	response, err := client.LookupFunction("0x8b477adb")
	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.Equal(t, "transferFromByLegacy(address,address,address,uint256)", response.Name)
}
