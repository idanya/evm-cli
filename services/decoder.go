package services

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/idanya/evm-cli/clients/directory"
	"github.com/idanya/evm-cli/entities"
)

type Decoder struct {
	directoryClient directory.DirectoryClient
}

func NewDecoder(directoryClient directory.DirectoryClient) *Decoder {
	return &Decoder{directoryClient}
}

func (d *Decoder) DecodeContractCallData(context context.Context, callData string) (*entities.DecodeResult, error) {
	payload := common.FromHex(callData)
	methodFourBytes := fmt.Sprintf("0x%s", common.Bytes2Hex(payload[:4]))
	methodPayload := payload[4:]

	lookupFunction, err := d.directoryClient.LookupFunction(methodFourBytes)
	if err != nil {
		return nil, err
	}

	inArgs, unpacked, err := UnpackFromSelector(lookupFunction, methodPayload)

	if err != nil {
		return nil, err
	}

	result := &entities.DecodeResult{}
	result.Method = lookupFunction
	result.Hash = methodFourBytes
	result.Arguments = make([]*entities.DecodedArgument, len(unpacked))
	for i, arg := range unpacked {
		result.Arguments[i] = d.toDecodedArgument(inArgs[i], arg)
	}

	return result, nil
}

func (cs *Decoder) toDecodedArgument(argument abi.Argument, value interface{}) *entities.DecodedArgument {
	return &entities.DecodedArgument{Name: argument.Name, Value: value, Type: argument.Type.String()}
}
