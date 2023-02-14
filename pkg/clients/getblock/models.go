package getblock

import "errors"

const (
	mainNetJsonRpc = "2.0"
	mainNetId      = "getblock.io"

	mainNetURL = "https://eth.getblock.io/mainnet/"
)

// MainNet Methods.
const (
	blockNumber      = "eth_blockNumber"
	getBlockByNumber = "eth_getBlockByNumber"
)

// MainNet Errors.
var (
	ErrTooManyRequests = errors.New("too many requests")
)

type MainNetRequest struct {
	JsonRpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      string        `json:"id"`
}

func NewMainNetRequest(method string, params ...interface{}) MainNetRequest {
	return MainNetRequest{
		JsonRpc: mainNetJsonRpc,
		Method:  method,
		Params:  params,
		ID:      mainNetId,
	}
}

type BlockResponse struct {
	JsonRpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  string `json:"result"`
}

type BlockByNumberResponse struct {
	JsonRpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Result  Block  `json:"result"`
}

type Block struct {
	BaseFeePerGas    string        `json:"baseFeePerGas"`
	Difficulty       string        `json:"difficulty"`
	ExtraData        string        `json:"extraData"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Hash             string        `json:"hash"`
	LogsBloom        string        `json:"logsBloom"`
	Miner            string        `json:"miner"`
	MixHash          string        `json:"mixHash"`
	Nonce            string        `json:"nonce"`
	Number           string        `json:"number"`
	ParentHash       string        `json:"parentHash"`
	ReceiptsRoot     string        `json:"receiptsRoot"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	Size             string        `json:"size"`
	StateRoot        string        `json:"stateRoot"`
	Timestamp        string        `json:"timestamp"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	Transactions     []Transaction `json:"transactions"`
	TransactionsRoot string        `json:"transactionsRoot"`
	Uncles           []string      `json:"uncles"`
}

type Transaction struct {
	BlockHash            string        `json:"blockHash"`
	BlockNumber          string        `json:"blockNumber"`
	From                 string        `json:"from"`
	Gas                  string        `json:"gas"`
	GasPrice             string        `json:"gasPrice"`
	Hash                 string        `json:"hash"`
	Input                string        `json:"input"`
	Nonce                string        `json:"nonce"`
	To                   string        `json:"to"`
	TransactionIndex     string        `json:"transactionIndex"`
	Value                string        `json:"value"`
	Type                 string        `json:"type"`
	V                    string        `json:"v"`
	R                    string        `json:"r"`
	S                    string        `json:"s"`
	MaxFeePerGas         string        `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas string        `json:"maxPriorityFeePerGas,omitempty"`
	AccessList           []interface{} `json:"accessList,omitempty"`
	ChainID              string        `json:"chainId,omitempty"`
}
