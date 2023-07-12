package services

import "github.com/ethereum/go-ethereum/accounts/abi"

func ArgumentsFromSelector(funcSig string) (abi.Arguments, error) {
	selector, err := abi.ParseSelector(funcSig)
	if err != nil {
		return nil, err
	}

	inArgs := make(abi.Arguments, len(selector.Inputs))
	for i, arg := range selector.Inputs {
		t, err := abi.NewType(arg.Type, arg.InternalType, arg.Components)
		if err != nil {
			return nil, err
		}

		inArgs[i] = abi.Argument{Name: arg.Name, Type: t, Indexed: arg.Indexed}
	}

	return inArgs, nil

}

func UnpackFromSelector(funcSig string, methodPayload []byte) (abi.Arguments, []interface{}, error) {
	inArgs, err := ArgumentsFromSelector(funcSig)
	if err != nil {
		return nil, nil, err
	}
	unpacked, err := inArgs.Unpack(methodPayload)
	return inArgs, unpacked, err
}
