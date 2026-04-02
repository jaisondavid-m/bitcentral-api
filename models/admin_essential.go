package models

type User struct {
	UID            string `json:"uid"`
	Email          string `json:"email"`
	DisplayName    string `json:"displayName"`
	PhotoURL       string `json:"photoURL"`
	CreationTime   string `json:"creationTime"`
	LastSignInTime string `json:"lastSignInTime"`
}

type PSToken struct {
	Token string `json:"token"`
}