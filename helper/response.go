package helper

import (
	"gitlab.com/learn-micorservices/profile-service/model/domain"
	"gitlab.com/learn-micorservices/profile-service/model/web"
)

// Profile Responses
func ToProfileResponse(user domain.User) web.ProfileResponse {
	return web.ProfileResponse{
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		Phone:    user.Phone,
		Role: web.RoleResponse{
			ID:   user.RoleID,
			Name: user.RoleName,
		},
	}
}