package coinlib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func GetConfigParameters() Configs {

	jsonFile, err := os.Open(configfile)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var lvls Configs

	json.Unmarshal(byteValue, &lvls)

	return lvls
}

func GetConfigLevel(c *Configs, whichlevel uint) *Config {
	conf := c.Configs[whichlevel]
	return &conf
}

func DumpConfigsInfo(l *Configs) {
	num_of_levels := len(l.Configs)
	fmt.Println("how many levels : ", num_of_levels)
	for _, each := range l.Configs {
		fmt.Println("name :", each.Name)
		fmt.Println("id   :", each.Id)
		fmt.Println("graph max data : ", each.Graph_max_data_points)
		fmt.Println("num data to start to record in cache : ", each.Max_num_days_to_record_incache)
		fmt.Println("num data to start realtime data : ", each.Num_data_to_start_realtime_data)
		fmt.Println("num mas :", each.Num_ma_calc)
		for _, x := range each.Malist {
			fmt.Println("         ma period found ", x.Period)

		}
	}
}
func GetApiSecret() (key string, secret string) {

	jsonFile, err := os.Open(configfile)
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var info ApiInfo
	json.Unmarshal(byteValue, &info)
	fmt.Println("secrets info :")

	return info.Secret.Key, info.Secret.Apisecret
}
