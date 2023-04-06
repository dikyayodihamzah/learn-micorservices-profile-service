package web

type RoleResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ProfileResponse struct {
	Name      string       `json:"name"`
	Username  string       `json:"username"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Role      RoleResponse `json:"role"`
}

