package coinlib

import "fmt"

type Wallet struct {
	Bucket map[string]float64
	Btc    float64 //btc equivalent of wallet

}

func (w *Wallet) IsNil() bool {
	if len(w.Bucket) == 0 {
		return true
	}
	return false
}
func (w *Wallet) Create() {
	w.Bucket = make(map[string]float64)
	w.Btc = 0
}

func (w *Wallet) Insert(currency string, amount float64) {
	w.Bucket[currency] = amount
}
func (w *Wallet) Dump() {
	fmt.Println("non zero balance currencies:")
	for k, v := range w.Bucket {
		fmt.Printf(" %s :%f \n", k, v)
	}
	if w.Btc != 0 {
		fmt.Println("BTC equivalent :", w.Btc)
	}
}
func (w *Wallet) FindCoinInWallet(coin string) float64 {
	if val, ok := w.Bucket[coin]; ok {
		return val
	}
	return 0
}
