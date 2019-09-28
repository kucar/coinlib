package coinlib

import "fmt"
import "time"

const (
	BUYING = iota + 1
	SELLING
	NEUTRAL
)

type RUNNER struct {
	Activity  int
	Professor *PROFESSOR
	Webserver *WEBSERVER
}

func (r *RUNNER) Init() {

	if r.Professor == nil {
		r.Professor = new(PROFESSOR)
	}
	if r.Webserver == nil {
		r.Webserver = new(WEBSERVER)
	}
	r.Activity = NEUTRAL
	r.Webserver.Init()
	r.Professor.Init(coinname, r.Webserver, interval)

}
func (r *RUNNER) Activator(newactivity int) {
	switch newactivity {
	case BUYING:
		fmt.Println("BUYING")
		r.Activity = BUYING
	case SELLING:
		fmt.Println("SELLING")
		r.Activity = SELLING
	case NEUTRAL:
		fmt.Println("NEUTRAL")
		r.Activity = NEUTRAL
	default:
		panic("unexpected activity type...")
	}
}
func (r *RUNNER) IsActivitySet() bool {
	switch r.Activity {
	case BUYING, SELLING:
		return true
	}
	return false
}
func (r *RUNNER) Run() {
	fmt.Println("main buyer/seller routine starts ")
	for {
		if false == r.IsActivitySet() {
			val := r.Professor.FindCoinInWallet(coinname)
			if val > 0 {
				//we have the coin in wallet, so look for selling option
				r.Activator(SELLING)
				//TODO initiate goroutine seller
			} else {
				r.Activator(BUYING)
				//TODO initiate goroutine buyer
			}
		} else {
			fmt.Println("waiting ", sleep_time, " seconds")
			time.Sleep(sleep_time * time.Second)
		}
	}
}
