package dto

type CreateUser struct {
	Name     string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateUser struct {
	Id       string `json:"id"`
	Password string `json:"password" binding:"required"`
}
