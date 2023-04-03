package web

import (
	"time"

	"gitlab.com/learn-micorservices/profile-service/model/domain"
)

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProfileResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Role      RoleResponse `json:"role"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func NewProfileResponse(profile domain.Profile) ProfileResponse {
	return ProfileResponse{
		ID:       profile.ID,
		Name:     profile.Name,
		Username: profile.Username,
		Email:    profile.Email,
		Phone:    profile.Phone,
		Role: RoleResponse{
			ID:   profile.Role.ID,
			Name: profile.Role.Name,
		},
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}
}
