package entities

type ERCTokenStandard struct {
	Name           string   `json:"name"`
	FunctionList   []string `json:"functionList"`
	FunctionHashes []string
}

type TokenStandards = []*ERCTokenStandard
