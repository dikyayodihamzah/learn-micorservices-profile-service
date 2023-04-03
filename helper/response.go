package helper

import (
	"gitlab.com/learn-micorservices/profile-service/model/domain"
	"gitlab.com/learn-micorservices/profile-service/model/web"
)

// Profile Responses
func ToProfileResponse(user domain.Profile) web.ProfileResponse {
	return web.ProfileResponse{
		ID:       user.ID,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: web.RoleResponse{
			ID:   user.Role.ID,
			Name: user.Role.Name,
		},
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}