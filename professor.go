package coinlib

import (
	"fmt"
	"math"
	"sync"
	"time"

	bittrex "github.com/toorop/go-bittrex"
)

type PROFESSOR struct {
	//interfaces
	client_intf    *BTRX
	webserver_intf *WEBSERVER
	ma_container   *MA_LIST
	//attributes
	coin_name       string
	historic_datasz int

	rt_price_data       []float64
	all_time_low        float64
	all_time_high       float64
	historic_price_data []float64
	interval_typ        string

	sync.RWMutex
}

func (p *PROFESSOR) Init(coinname string, webserver *WEBSERVER, interval_typ string) {
	DumpConfigsInfo(&configs)
	p.coin_name = coinname
	p.client_intf = new(BTRX)
	p.client_intf.Init()
	p.webserver_intf = webserver
	p.interval_typ = interval_typ
	p.all_time_high = math.MaxFloat64 * -1
	p.all_time_low = math.MaxFloat64
	p.InitPriceList()
	p.historic_datasz = len(p.historic_price_data)
	println("Professor::Init> size of historic data set ", p.historic_datasz)
	p.initMAContainer()
	ticker := time.NewTicker(time.Duration(UPDATE_TIME_INTERVAL_SEC) * time.Second)
	go func() {
		for range ticker.C {
			p.GetLatestTickPeriodic()
			p.LatestMovAvgPeriodic()
			p.Monitor()
			p.UpdateWebServer()
		}
	}()

}
func (p *PROFESSOR) GetAvgABL(t bittrex.Ticker) float64 {
	x, _ := t.Ask.Float64()
	y, _ := t.Bid.Float64()
	z, _ := t.Last.Float64()
	return (x + y + z) / 3
}
func (p *PROFESSOR) GetAvgHLC(t bittrex.Candle) float64 {
	x, _ := t.High.Float64()
	y, _ := t.Low.Float64()
	z, _ := t.Close.Float64()
	return (x + y + z) / 3
}
func (p *PROFESSOR) GetLatestTickPeriodic() {
	p.RLock()
	defer p.RUnlock()
	// its a tick update
	marketname := "BTC-" + coinname
	t, _ := p.client_intf.btrx_ptr.GetTicker(marketname)
	if len(p.rt_price_data) == MAX_PRICE_DATASZ {
		p.rt_price_data = p.rt_price_data[1:]
	}

	latest_tick := p.GetAvgABL(t)
	fmt.Println("Professor::GetLatestTickPeriodic> latest tick received :", latest_tick)
	p.rt_price_data = append(p.rt_price_data, latest_tick)

	if len(p.rt_price_data) > MAX_PRICE_DATASZ {
		panic("size of price data slice increasing")
	}
	if latest_tick > p.all_time_high {
		p.all_time_high = latest_tick
	} else if latest_tick < p.all_time_low {
		p.all_time_low = latest_tick
	}

}

func (p *PROFESSOR) LatestMovAvgPeriodic() {
	p.RLock()
	defer p.RUnlock()
	latest_tick := p.rt_price_data[len(p.rt_price_data)-1]

	for _, sma := range p.ma_container.ma_lst {
		//just to prevent grow to infinity
		if len(sma.ma_data) == MAX_MOVAVG_DATASZ {
			sma.ma_data = sma.ma_data[1:]
		}
		sma.Add(latest_tick)  //add the tick
		sma.Append(sma.Avg()) //add the latest ma
		if len(sma.ma_data) > MAX_MOVAVG_DATASZ {
			panic("size of ma data  slice increasing")
		}
	}
}
func (p *PROFESSOR) Monitor() {
	fmt.Println("Professor::Monitor> price data len = ", len(p.rt_price_data))
	fmt.Println("Professor::Monitor> last price data = ", p.rt_price_data[len(p.rt_price_data)-1], " min= ", p.all_time_low, " max= ", p.all_time_high)
	masize := len(p.ma_container.ma_lst)
	fmt.Printf("Professor::Monitor> ma container size = %d\n", masize)
	for _, eachma := range p.ma_container.ma_lst {
		ma_max, _ := eachma.Max()
		fmt.Println("Professor::Monitor> last ma_", eachma.window_size_str, " =", eachma.ma_data[len(eachma.ma_data)-1], "max= ", ma_max)
	}

}

func (p *PROFESSOR) UpdateWebServer() {
	p.webserver_intf.ticker_chan <- p.rt_price_data
	fmt.Println("sent ticker to webserver")
	p.webserver_intf.ma_chan <- p.ma_container
	fmt.Println("sent mov_avg to webserver")

}
func (p *PROFESSOR) InitPriceList() {
	marketname := "BTC-" + coinname
	if len(p.historic_price_data) != 0 { // its an init
		panic("historic data expected to be nil")
	}

	candles, err := p.client_intf.btrx_ptr.GetTicks(marketname, p.interval_typ)
	if err != nil {
		fmt.Println(err)
		panic("historic price data could not be gathered")
	}
	if candles == nil {
		panic("candles empty after GetTicks()")
	}
	for _, candle := range candles {
		price := p.GetAvgHLC(candle)
		if price > p.all_time_high {
			p.all_time_high = price
		} else if price < p.all_time_low {
			p.all_time_low = price
		}
		p.historic_price_data = append(p.historic_price_data, price)
	}
	//trick: we need to get latest MA size (eg 50) points from historic
	//        data to real time to make the mov avg calculations
	//        satisfied in the beginning
	//deep cpy historic to real_time
	fmt.Println("historic data initialized , len ", len(p.historic_price_data))
	p.rt_price_data = append(p.rt_price_data[:0:0], p.historic_price_data...)
	p.rt_price_data = p.rt_price_data[len(p.rt_price_data)-NUMBER_DATA_TO_START_RT_SERIES : len(p.rt_price_data)]
	fmt.Println("len of rt_price_data ", len(p.rt_price_data))
	p.webserver_intf.ticker_chan <- p.rt_price_data

}

func (p *PROFESSOR) initMAContainer() {

	if len(p.rt_price_data) == 0 {
		panic("price data list expected to be empty")
	}
	p.ma_container = new(MA_LIST)

	ma_list_to_be_created := p.GenerateMaSpliceFromConfig(config)
	for _, val := range ma_list_to_be_created {
		sma := new(MA_SIMPLE)
		sma.Init(val)
		for _, price := range p.rt_price_data[NUMBER_DATA_TO_START_RT_SERIES-sma.window_size_int:] {
			sma.Add(price)
		}
		p.ma_container.Insert(sma)
	}

}

func (p *PROFESSOR) FindCoinInWallet(str string) float64 {
	return p.client_intf.wallet.FindCoinInWallet(str)
}

func (p *PROFESSOR) GenerateMaSpliceFromConfig(lb *Config) []int {
	ret := make([]int, 0)
	for _, ma := range lb.Malist {
		ret = append(ret, ma.Period)
	}
	return ret
}
