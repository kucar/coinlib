package coinlib

type Config struct {
	Name                            string `json:"name"`
	Id                              int    `json:"id"`
	Graph_max_data_points           int    `json:"graph_max_data_points"`
	Num_data_to_start_realtime_data int    `json:"num_data_to_start_realtime_data"`
	Max_num_days_to_record_incache  int    `json:"max_num_days_to_record_incache"`
	Update_time_interval_sec        int    `json:"update_time_interval_sec"`
	Num_ma_calc                     int    `json:"num_ma_calc"`
	Malist                          []ma   `json:"malist"`
}

//simplex ma type ie : 30 days MA
type ma struct {
	Period int `json:"period"`
}
