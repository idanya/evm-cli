package openchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Client struct {
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) LookupFunction(hash string) (*LookupFunction, error) {
	requestUrl := fmt.Sprintf("https://api.openchain.xyz/signature-database/v1/lookup?function=%s&filter=true", hash)
	url, err := url.Parse(requestUrl)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result LookupResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	if result.Ok && len(result.Result.Function[hash]) > 0 {
		return result.Result.Function[hash][0], nil
	}

	return nil, fmt.Errorf("Function not found")
}
