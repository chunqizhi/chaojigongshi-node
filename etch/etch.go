package etch

import (
	"github.com/chunqizhi/go-ethereum/core/types"
	"github.com/chunqizhi/go-ethereum/ethclient"
	"golang.org/x/net/context"
	"math/big"
)

type Eclient struct {
	*ethclient.Client
}

func New(url string) (*Eclient, error) {
	this := new(Eclient)
	client, err := ethclient.Dial(url)
	if err != nil {
		return this, err
	}
	this.Client = client
	return this, nil
}
func (this Eclient) Count() (int64, error) {
	header, err := this.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return 0, err
	}
	return header.Number.Int64(), nil
}
func (this Eclient) Block(number int64) (*types.Block, error) {
	//ceshi := header.Number.Int64()
	blockNumber := big.NewInt(number)
	block, err := this.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		return block, err
	}
	return block, nil
}

func (this Eclient) GetTransactionReceipt(hash string) (*types.Receipt, error) {
	return this.TransactionReceipt(context.Background(), hash)
}
