package etherscan

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"gitlab.com/fireblocks/web3/utils/evm-cli/entities"
)

type GenericEtherscanClient struct {
	baseUrl string
	apiKey  string
}

func NewGenericEtherscanClient(baseUrl string, apiKey string) *GenericEtherscanClient {
	return &GenericEtherscanClient{baseUrl, apiKey}
}

func (c *GenericEtherscanClient) sendRequest(context context.Context, module string, action string, extraParams map[string]string) ([]byte, error) {
	url, err := url.Parse(c.baseUrl)
	if err != nil {
		return nil, err
	}

	queryValue := url.Query()
	queryValue.Set("module", module)
	queryValue.Set("action", action)
	queryValue.Set("apikey", c.apiKey)

	for k, v := range extraParams {
		queryValue.Set(k, v)
	}

	url.RawQuery = queryValue.Encode()

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *GenericEtherscanClient) GetContractABI(context context.Context, contractAddress string) (entities.ABI, error) {
	log.Printf("Fetching ABI for contract address %s", contractAddress)
	data, err := c.sendRequest(context, "contract", "getabi", map[string]string{"address": contractAddress})
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	var abi entities.ABI
	if err := json.Unmarshal([]byte(result["result"].(string)), &abi); err != nil {
		return nil, err
	}

	return abi, nil
}
