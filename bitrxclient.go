package coinlib

import (
	"fmt"
	"time"

	bittrex "github.com/toorop/go-bittrex"
)

const comission float64 = 0.0025

type BTRX struct {
	btrx_ptr *bittrex.Bittrex
	wallet   Wallet
}

func (b *BTRX) Updatewallet() {
	if b.wallet.IsNil() == true {
		fmt.Println("wallet is nill,creating..")
		b.wallet.Create()
	}
	balances, _ := b.btrx_ptr.GetBalances()
	for _, balance := range balances {
		x, _ := balance.Balance.Float64()
		if x != 0 {
			b.wallet.Insert(balance.Currency, x)
		}

	}
	b.wallet.Btc = b.Calculate_btc()
}
func (b *BTRX) Init() {
	// Bittrex client
	if b.btrx_ptr == nil {
		k, s := GetApiSecret()
		b.btrx_ptr = bittrex.New(k, s)
	}
	b.Updatewallet()
	b.wallet.Dump()
}
func (b *BTRX) Calculate_btc() float64 {
	var sum float64 = 0
	for k, v := range b.wallet.Bucket {
		ticker, _ := b.btrx_ptr.GetTicker("BTC-" + k)
		ticker_val, _ := ticker.Last.Float64()
		sum += v * ticker_val
	}
	return sum
}
func (b *BTRX) Monitor_gr() {

	ticker := time.NewTicker(5 * time.Second)

	for t := range ticker.C {
		fmt.Println("Monitor_gr triggered", t)
		b.Updatewallet()
		b.wallet.Dump()
	}

}
