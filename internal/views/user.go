package views

type User struct {
	Id      uint   `json:"id"`
	Name    string `json:"username"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"isAdmin"`
}
