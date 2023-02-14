package getblock

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
)

type Mainnet interface {
	GetBlockNumber(ctx context.Context) (BlockResponse, error)
	GetBlockByNumber(ctx context.Context, address int64) (*BlockByNumberResponse, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	apiKey     string
	httpClient HTTPClient
}

var _ Mainnet = Client{}

func NewClient(apiKey string, httpClient HTTPClient) *Client {
	return &Client{
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (c Client) GetBlockNumber(ctx context.Context) (BlockResponse, error) {
	request := NewMainNetRequest(blockNumber)

	requestData, err := json.Marshal(request)
	if err != nil {
		return BlockResponse{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mainNetURL, bytes.NewReader(requestData))
	if err != nil {
		return BlockResponse{}, err
	}

	req.Header.Add("x-api-key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return BlockResponse{}, err
	}
	defer res.Body.Close()

	var result BlockResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return BlockResponse{}, err
	}

	return result, nil
}

func (c Client) GetBlockByNumber(ctx context.Context, address int64) (*BlockByNumberResponse, error) {
	addressHex := "0x" + strconv.FormatInt(address, 16)

	request := NewMainNetRequest(getBlockByNumber, addressHex, true)

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, mainNetURL, bytes.NewReader(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("x-api-key", c.apiKey)
	req.Header.Add("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		return nil, ErrTooManyRequests
	}

	var result BlockByNumberResponse
	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
