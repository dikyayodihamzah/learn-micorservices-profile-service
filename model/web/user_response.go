package web

import (
	"time"

	"gitlab.com/learn-micorservices/user-service/model/domain"
)

type UserRoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserResponse struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Username  string           `json:"username"`
	Email     string           `json:"email"`
	Phone     string           `json:"phone"`
	Role      UserRoleResponse `json:"role"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func NewuserResponse(user domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: UserRoleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
