package entities

import (
	"fmt"
	"strings"
)

type ABIParameter struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (p *ABIParameter) String() string {
	return p.Name + " " + p.Type
}

type ABIParameterList []ABIParameter

func (l ABIParameterList) String() string {
	types := make([]string, len(l))
	for i, input := range l {
		types[i] = input.String()
	}

	return strings.Join(types, ", ")
}

type ABIFunction struct {
	Type            string           `json:"type"`
	Name            string           `json:"name"`
	Payable         bool             `json:"payable"`
	StateMutability string           `json:"stateMutability"`
	Constant        bool             `json:"constant"`
	Inputs          ABIParameterList `json:"inputs,omitempty"`
	Outputs         ABIParameterList `json:"outputs,omitempty"`
}

func (f *ABIFunction) String() string {
	return fmt.Sprintf("[%s] %s(%s) returns (%s)", f.StateMutability, f.Name, f.Inputs.String(), f.Outputs.String())
}

type ABI []ABIFunction
