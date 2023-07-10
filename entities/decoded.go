package entities

type DecodedArgument struct {
	Name       string
	Value      interface{}
	Type       string		
}

type DecodeResult struct {
	Method    string
	Hash      string
	Arguments []*DecodedArgument
}
