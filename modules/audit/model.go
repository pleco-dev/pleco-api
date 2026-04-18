package audit

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	ActorUserID *uint  `json:"actor_user_id"`
	Action      string `json:"action"`
	Resource    string `json:"resource"`
	ResourceID  *uint  `json:"resource_id"`
	Status      string `json:"status"`
	Description string `json:"description"`
	IPAddress   string `json:"ip_address"`
	UserAgent   string `json:"user_agent"`
}
