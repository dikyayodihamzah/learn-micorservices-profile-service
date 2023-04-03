package web

type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

type UpdateUsernameRequest struct {
	Username string `json:"username"`
}

type UpdateRoleRequest struct {
	RoleID   string `json:"role_id"`
}
