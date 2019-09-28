package coinlib

const SECONDS_IN_DAY int = 86400

var configs Configs = GetConfigParameters()
var config *Config = GetConfigLevel(&configs, 0)
var NUM_MA_CALC int = config.Num_ma_calc
var MAX_NUM_DAYS_TO_RECORD int = config.Max_num_days_to_record_incache
var UPDATE_TIME_INTERVAL_SEC int = config.Update_time_interval_sec
var NUMBER_DATA_TO_START_RT_SERIES int = config.Num_data_to_start_realtime_data
var MAX_NUM_DATA_PTS_IN_CACHE int = (MAX_NUM_DAYS_TO_RECORD * SECONDS_IN_DAY) / UPDATE_TIME_INTERVAL_SEC
var MAX_MOVAVG_DATASZ int = MAX_NUM_DATA_PTS_IN_CACHE
var MAX_PRICE_DATASZ int = MAX_NUM_DATA_PTS_IN_CACHE

const coinname string = "BCH"

// Interval can be -> ["oneMin", "fiveMin", "thirtyMin", "hour", "day"]
const interval string = "oneMin"
const sleep_time = 5
