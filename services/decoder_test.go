package services

import (
	"context"
	"encoding/json"
	"log"
	"testing"

	dirmock "github.com/idanya/evm-cli/clients/directory/mocks"
	"github.com/stretchr/testify/assert"
)

func TestContractService_DecodeContractCallData(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)

	directoryClientMock.On("LookupFunction", "0x791ac947").Return("swapExactTokensForETHSupportingFeeOnTransferTokens(uint256,uint256,address[],address,uint256)", nil)

	callData := "0x791ac9470000000000000000000000000000000000000000002fb44e7672f91c317f000000000000000000000000000000000000000000000000000000bf85a503ff565100000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000300e6122f18d52a5052c5227dad5204669f37cfd0000000000000000000000000000000000000000000000000000000064abc05000000000000000000000000000000000000000000000000000000000000000020000000000000000000000000981d9774a59a703db85f5eaa23672283ea31106000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"

	decoded, err := decoder.DecodeContractCallData(context.Background(), callData)
	assert.NoError(t, err)

	assert.Equal(t, "swapExactTokensForETHSupportingFeeOnTransferTokens(uint256,uint256,address[],address,uint256)", decoded.Method)
	assert.Equal(t, "0x791ac947", decoded.Hash)
}

func TestContractService_DecodeContractCallData_Tuples(t *testing.T) {
	directoryClientMock := dirmock.NewDirectoryClient(t)
	decoder := NewDecoder(directoryClientMock)
	directoryClientMock.On("LookupFunction", "0xe7acab24").Return("fulfillAdvancedOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes),(uint256,uint8,uint256,uint256,bytes32[])[],bytes32,address)", nil)

	callData := "0xe7acab24000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000005e00000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000046000000000000000000000000000000000000000000000000000000000000004c00000000000000000000000004a1d593921d76e8ff52f2f7026e3bf1765b0722e000000000000000000000000000000e7ec00e7b300774b00001314b8610022b80000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000022000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000064ab457e0000000000000000000000000000000000000000000000000000000064af39e30000000000000000000000000000000000000000000000000000000000000000360c6ebe00000000000000000000000000000000000000001711c345c46c535f0000007b02230091a7ed01230072f7006a004d60a8d4e71d599b8104250f0000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000058d15e176280000000000000000000000000000000000000000000000000000058d15e17628000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030000000000000000000000000e9e3eff16d2d19f5d7a04717d70b81f2c7465a500000000000000000000000000000000000000000000000000000000000001a5000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000004a1d593921d76e8ff52f2f7026e3bf1765b0722e0000000000000000000000000000000000000000000000000000000000000001000000000000000000000000c02aaa39b223fe8d0a0e5c4f27ead9083c756cc200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000002386f26fc10000000000000000000000000000000000000000000000000000002386f26fc10000000000000000000000000000000a26b00c1f0df003000390027140000faa719000000000000000000000000000000000000000000000000000000000000004028432de29d4bc1e7313ccf9b6dba2296ed9716ff072ec8be51ecc0a0464499a6e2837ad2e0d591dcb3254d23de18f2827e923d174bd5bac02241a2a86af69acb000000000000000000000000000000000000000000000000000000000000007e004c67d20778f55b8f7078a6c59cb69d6e610dadb80000000064abced9678f61c4f78c2fe474d99df7632479f14e8ab1afb991db14e78857951c00fc598fb3929a85b2667b6903d76b76b9112aeaa309003b380a793c1dc07218b839180000000000000000000000000000000000000000000000000000000000000001a50000000000000000000000000000000000000000000000000000000000000000000000000000360c6ebe"

	decoded, err := decoder.DecodeContractCallData(context.Background(), callData)
	assert.NoError(t, err)

	decodedJson, err := json.MarshalIndent(decoded, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Decoded call data:\n\n%s", decodedJson)

	assert.Equal(t, "fulfillAdvancedOrder(((address,address,(uint8,address,uint256,uint256,uint256)[],(uint8,address,uint256,uint256,uint256,address)[],uint8,uint256,uint256,bytes32,uint256,bytes32,uint256),uint120,uint120,bytes,bytes),(uint256,uint8,uint256,uint256,bytes32[])[],bytes32,address)", decoded.Method)
	assert.Equal(t, "0xe7acab24", decoded.Hash)
}