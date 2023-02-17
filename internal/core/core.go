package core

import "math/big"

type Transaction struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"amount"`
}
