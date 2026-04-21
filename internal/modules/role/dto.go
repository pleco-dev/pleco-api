package role

type UpdateRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

type RolePermissionsResponse struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}
