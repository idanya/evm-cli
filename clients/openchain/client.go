package openchain

type LookupFunction struct {
	Name     string `json:"name"`
	Filtered bool   `json:"filtered"`
}

type LookupResponseResult struct {
	Function map[string][]*LookupFunction `json:"function"`
}

type LookupResponse struct {
	Ok     bool                 `json:"ok"`
	Result LookupResponseResult `json:"result"`
}

type OpenChainClient interface {
	LookupFunction(hash string) (*LookupFunction, error)
}
