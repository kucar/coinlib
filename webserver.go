package coinlib

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"sync"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

var inc int = 0
var inc_ma int = 0

const GRAPH_DATA_POINTS_MAX int = 450 //todo config

var ma_color_map = map[int]color.RGBA{
	30:  color.RGBA{0x00, 0x80, 0x00, 0xff},
	50:  color.RGBA{0xdb, 0x70, 0x93, 0xff},
	80:  color.RGBA{0x99, 0x32, 0xcc, 0xff},
	100: color.RGBA{0x99, 0x00, 0xcc, 0xff},
	150: color.RGBA{0x99, 0x20, 0xcc, 0xff},
	300: color.RGBA{0x99, 0x24, 0xab, 0xba},
	350: color.RGBA{0x92, 0x45, 0x23, 0x11},
	600: color.RGBA{0x09, 0xca, 0xac, 0x1f}}

type WEBSERVER struct {
	ticker_data []float64
	sync.RWMutex
	ticker_chan  chan []float64 //channel of slice
	ma_chan      chan *MA_LIST  //channel of ma container (ptr)
	ma_container *MA_LIST
}

func (s *WEBSERVER) Init() {
	s.ticker_chan = make(chan []float64, 1)
	s.ma_chan = make(chan *MA_LIST, 1)
	s.datacollect()
	http.HandleFunc("/", s.root)
	http.HandleFunc("/statz", s.statz)
	http.HandleFunc("/statz/cizici.png", s.scatter)
	s.serve()
}
func (s *WEBSERVER) root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "--- webserver ---")

}

func (s *WEBSERVER) scatter(w http.ResponseWriter, r *http.Request) {
	s.RLock()
	defer s.RUnlock()
	/*
		Main Plot
	*/
	p, err := plot.New()
	p.Title.Text = "BTC-" + coinname
	p.Y.Label.Text = "btc "
	p.X.Label.Text = "time series "
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	/*
		TICKER GRAPH
	*/

	if len(s.ticker_data) > GRAPH_DATA_POINTS_MAX {
		inc = len(s.ticker_data) - GRAPH_DATA_POINTS_MAX
	}

	graph_ticker := s.ticker_data[inc:]
	xys := make(plotter.XYs, len(graph_ticker))

	for i, d := range graph_ticker {
		xys[i].X = float64(i)
		xys[i].Y = float64(d)

	}
	sc, err := plotter.NewScatter(xys)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sc.GlyphStyle.Shape = draw.CrossGlyph{}
	p.Add(sc)
	/*
		MA GRAPH
	*/
	if nil == s.ma_container {
		fmt.Println("WEBSERVER:: data not ready yet ..")
		return
	}

	for _, sma := range s.ma_container.ma_lst {
		var offset float64 = float64(NUMBER_DATA_TO_START_RT_SERIES) - float64(inc)
		if len(sma.ma_data) > GRAPH_DATA_POINTS_MAX {
			inc_ma = len(sma.ma_data) - GRAPH_DATA_POINTS_MAX
			if offset < 0 {
				offset = 0
			}
		}
		ma_temp := sma.ma_data[inc_ma:]
		mas := make(plotter.XYs, len(ma_temp))

		for i, d := range ma_temp {
			mas[i].X = float64(i) + float64(offset)
			mas[i].Y = float64(d)
		}

		ma_line, err := plotter.NewLine(mas)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ma_line.Color = sma.renk
		p.Add(ma_line)
		p.Legend.Add(sma.window_size_str, ma_line)
		// ind += 1
	}

	wt, err := p.WriterTo(512, 512, "png")
	_ = wt
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	wt.WriteTo(w)
	for _, d := range s.ticker_data {
		fmt.Fprintln(w, d)
	}
	// p.Save(4*vg.Inch, 4*vg.Inch, "gocoin/scatter.png")
}

func (s *WEBSERVER) datacollect() {
	go func() {
		for {
			select {
			case tick := <-s.ticker_chan:
				s.ticker_data = tick
			case ma := <-s.ma_chan:
				s.ma_container = ma

			}
		}
	}()
}

func (s *WEBSERVER) serve() {
	go func() {
		log.Fatal(http.ListenAndServe(":9090", nil))
	}()
}

func (s *WEBSERVER) statz(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", `
<h1> webserver : </h1>
<h3> auto refresh rate 1 sec </h3>
<img src="/statz/cizici.png?rand=0" style="width:50%">
	<script>
setInterval(function() {
	var imgs = document.getElementsByTagName("IMG");
	for( var i = 0 ; i < imgs.length ; i++){
		var eqPos = imgs[i].src.lastIndexOf("=");
		var src   = imgs[i].src.substr(0, eqPos+1);
		imgs[i].src = src + Math.random();
	}
},1000);
	</script>	`)
}

//gocoin/scatter.png
