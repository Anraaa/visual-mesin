package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	LastPage    int   `json:"last_page"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  UserDTO `json:"user"`
}

type UserDTO struct {
	ID        uint    `json:"id"`
	NIP       *string `json:"nip"`
	UserName  string  `json:"user_name"`
	UserLevel string  `json:"user_level"`
	Email     *string `json:"email"`
	Roles     []Role  `json:"roles,omitempty"`
}

type RegisterRequest struct {
	UserName string `json:"user_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	UserLevel string `json:"user_level" binding:"omitempty,oneof=admin eng tech prod"`
}
