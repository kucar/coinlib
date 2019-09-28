package coinlib

import (
	"fmt"
	"image/color"

	"strconv"

	malib "github.com/RobinUS2/golang-moving-average"
)

type MA_SIMPLE struct {
	window_size_int int
	window_size_str string
	ma              *malib.MovingAverage
	ma_data         []float64
	renk            color.RGBA
}

func (s *MA_SIMPLE) Init(window_sz int) {
	s.window_size_int = window_sz
	s.renk = ma_color_map[window_sz]
	s.window_size_str = strconv.Itoa(window_sz)
	fmt.Println("MA_SIMPLE::Init() moving average with window size ", s.window_size_str, " is being initiated ...")
	s.ma = malib.New(window_sz)
	s.ma_data = make([]float64, 0, MAX_MOVAVG_DATASZ)
}
func (s *MA_SIMPLE) Add(values ...float64) {
	s.ma.Add(values...)
}
func (s *MA_SIMPLE) Avg() float64 {
	return s.ma.Avg()
}
func (s *MA_SIMPLE) Append(values ...float64) int {
	s.ma_data = append(s.ma_data, values...)
	return len(s.ma_data)
}
func (s *MA_SIMPLE) Max() (float64, error) {
	return s.ma.Max()
}
