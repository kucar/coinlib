package coinlib

type ApiInfo struct {
	Secret secret `json:"secret"`
}
type secret struct {
	Key       string `json:"apikey"`
	Apisecret string `json:"apisecret"`
}
