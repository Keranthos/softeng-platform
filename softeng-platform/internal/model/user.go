package model

import (
	"time"
)

type User struct {
	ID          int       `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	Nickname    string    `json:"nickname" db:"nickname"`
	Email       string    `json:"email" db:"email"`
	Password    string    `json:"-" db:"password"`
	Avatar      string    `json:"avater" db:"avatar"`
	Description string    `json:"description" db:"description"`
	FacePhoto   string    `json:"face_photo" db:"face_photo"`
	Role        string    `json:"role" db:"role"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type RegisterRequest struct {
	Username        string `form:"username" json:"username" binding:"required"`
	Email           string `form:"email" json:"email" binding:"required,email"`
	Password        string `form:"password" json:"password" binding:"required,min=6"`
	EmailPassword   string `form:"email_password" json:"email_password" binding:"required"`
	CertifyPassword string `form:"certify_password" json:"certify_password" binding:"required"`
}

type LoginRequest struct {
	UsernameOrEmail string `form:"username_or_email" json:"username_or_email" binding:"required"`
	Password        string `form:"password" json:"password" binding:"required"`
}

type UpdateProfileRequest struct {
	Nickname    string `form:"nickname" json:"nickname"`
	Avatar      string `form:"avater" json:"avater"`
	Description string `form:"description" json:"description"`
	FacePhoto   string `form:"face_photo" json:"face_photo"`
}
