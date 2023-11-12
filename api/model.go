package api

type User struct {
	Name      string   `json:"name"`
	Phones    []string `json:"phones"`
	Addresses []string `json:"addresses"`
}
