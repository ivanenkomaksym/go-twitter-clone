package models

type User struct {
	IsAnonymous bool   `json:"-"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Picture     string `json:"picture"`
}
