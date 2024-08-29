package dto

import (
	"time"

	"github.com/tricong1998/go-ecom/internal/models"
)

type CreateUserDto struct {
	Username string `json:"username" binding:"required"`
	FullName string `json:"full_name" binding:"required"`
}

type ReadUserRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	FullName  string    `json:"full_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListUserQuery struct {
	Username *string `form:"username"`
	Page     int32   `form:"page" binding:"required,min=1"`
	PerPage  int32   `form:"per_page" binding:"required,min=5,max=10"`
}

type ListUserResponse struct {
	Items    []UserResponse `json:"items"`
	Metadata MetadataDto    `json:"metadata"`
}

func ToUserResponse(user *models.User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
