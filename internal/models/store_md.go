package models

type Store struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Balans   int    `json:"Balans"`
}
