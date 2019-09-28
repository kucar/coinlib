package coinlib

const ERROR_MA_NOTFOUND int = 900
const ERROR_MA_FOUND int = 200

type MA_LIST struct {
	ma_lst   map[int]*MA_SIMPLE
	lst_size int
}

func (m *MA_LIST) Insert(simplex_ma *MA_SIMPLE) int {
	if nil == simplex_ma {
		panic("MA_LIST::Insert() pointer is NULL")
	}
	if nil == m.ma_lst {
		m.ma_lst = make(map[int]*MA_SIMPLE)
	}
	m.ma_lst[simplex_ma.window_size_int] = simplex_ma
	m.lst_size = len(m.ma_lst)
	return m.lst_size
}
func (m *MA_LIST) Find_ma(window_sz int) (*MA_SIMPLE, int) {
	if val, ok := m.ma_lst[window_sz]; ok {
		return val, ERROR_MA_FOUND
	}
	return nil, ERROR_MA_NOTFOUND
}
