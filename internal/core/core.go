package core

import "math/big"

type Transaction struct {
	From        string   `json:"from"`
	To          string   `json:"to"`
	AmountInWei *big.Int `json:"amount"`
}

type Error struct {
	Message string `json:"message"`
}
