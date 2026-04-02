package models

type Card struct {
	Image    string   `json:"img"`
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
	Link     string   `json:"link"`
	BtnText  string   `json:"btntext"`
}