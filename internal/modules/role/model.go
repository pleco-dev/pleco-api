package role

type Role struct {
	ID              uint
	Name            string
	RolePermissions []RolePermission `gorm:"foreignKey:RoleID" json:"role_permissions,omitempty"`
}

type RolePermission struct {
	ID         uint   `json:"id"`
	RoleID     uint   `json:"role_id"`
	Permission string `json:"permission"`
}
