package coinlib

var configfile string = "./config/config.json"

//config container
type Configs struct {
	Configs []Config `json:"configs"`
}
