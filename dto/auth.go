package dto

type RegisterRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	UserName string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
