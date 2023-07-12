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

func (c *Client) LookupFunction(hash string) (string, error) {
	return c.lookupEntity(EntityTypeFunction, hash)
}

func (c *Client) LookupEvent(hash string) (string, error) {
	return c.lookupEntity(EntityTypeEvent, hash)
}

func (c *Client) lookupEntity(entityType EntityType, hash string) (string, error) {
	requestUrl := fmt.Sprintf("https://api.openchain.xyz/signature-database/v1/lookup?%s=%s&filter=true", entityType, hash)
	url, err := url.Parse(requestUrl)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result LookupResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	if result.Ok {
		switch entityType {
		case EntityTypeFunction:
			if len(result.Result.Function[hash]) > 0 {
				return result.Result.Function[hash][0].Name, nil
			}
		case EntityTypeEvent:
			if len(result.Result.Event[hash]) > 0 {
				return result.Result.Event[hash][0].Name, nil
			}
		}
	}

	return "", fmt.Errorf("Entity not found")
}
