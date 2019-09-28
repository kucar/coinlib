package main

import (
	_ "net/http/pprof"
	"sync"

	"github.com/kucar/coinlib"
)

var wg sync.WaitGroup

func main() {

	wg.Add(1)
	runner := new(coinlib.RUNNER)
	runner.Init()
	runner.Run()

	wg.Wait()
}
