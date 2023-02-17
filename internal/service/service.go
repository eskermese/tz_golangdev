package service

import (
	"context"
	"errors"
	"math/big"
	"strconv"
	"strings"
	"sync"

	"github.com/eskermese/tz_golangdev/internal/core"
	"github.com/eskermese/tz_golangdev/pkg/clients/getblock"
	"github.com/eskermese/tz_golangdev/pkg/workers"
)

const NumBlocksToInclude = 100

type Deps struct {
	Getblock getblock.Mainnet
}

type Service struct {
	getblock getblock.Mainnet
}

func New(deps Deps) *Service {
	return &Service{
		getblock: deps.Getblock,
	}
}

func (s Service) GetMaxTransactionChange(ctx context.Context) (core.Transaction, error) {
	blockNumber, err := s.getblock.GetBlockNumber(ctx)
	if err != nil {
		return core.Transaction{}, err
	}

	blockNumberDec, err := strconv.ParseInt(hexaNumberToInteger(blockNumber.Result), 16, 32)
	if err != nil {
		return core.Transaction{}, err
	}

	firstBlockNumberDec := blockNumberDec - NumBlocksToInclude

	mu := &sync.Mutex{}

	// Run goroutines for getting transactions from blocks
	g, gCtx := workers.GroupWithContext(ctx)

	results := make([]core.Transaction, 0, NumBlocksToInclude)

	for i := firstBlockNumberDec; i < blockNumberDec; i++ {
		j := i

		g.Go(func() error {
			result, err := s.getMaxValueResultFromBlock(gCtx, j)
			if err != nil {
				// TODO: Чтобы обработат нормально, можно было бы использовать общие состояние,
				// такие как stopped/paused/running, чтобы при ErrTooManyRequests, изменять running на paused на определноое время.
				// И по истечению вернуть на running
				if errors.Is(err, getblock.ErrTooManyRequests) {
					return nil
				}

				return err
			}

			mu.Lock()
			results = append(results, result)
			mu.Unlock()

			return nil
		})
	}

	if err = g.Wait(); err != nil {
		return core.Transaction{}, err
	}

	result := core.Transaction{
		Amount: big.NewInt(0),
	}
	for _, v := range results {
		if v.Amount.Cmp(result.Amount) == 1 {
			result = v
		}
	}

	return result, nil
}

func (s Service) getMaxValueResultFromBlock(ctx context.Context, address int64) (core.Transaction, error) {
	blockData, err := s.getblock.GetBlockByNumber(ctx, address)
	if err != nil {
		return core.Transaction{}, err
	}

	result := core.Transaction{
		Amount: big.NewInt(0),
	}

	for _, transaction := range blockData.Result.Transactions {
		value := new(big.Int)
		value.SetString(hexaNumberToInteger(transaction.Value), 16)

		if value.Cmp(result.Amount) == 1 {
			result.Amount = value
			result.From = transaction.From
			result.To = transaction.To
		}
	}

	return result, nil
}

func hexaNumberToInteger(hexaString string) string {
	numberStr := strings.ReplaceAll(hexaString, "0x", "")
	numberStr = strings.ReplaceAll(numberStr, "0X", "")

	return numberStr
}
