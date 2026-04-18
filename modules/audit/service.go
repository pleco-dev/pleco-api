package audit

import "log"

type RecordInput struct {
	ActorUserID *uint
	Action      string
	Resource    string
	ResourceID  *uint
	Status      string
	Description string
	IPAddress   string
	UserAgent   string
}

type Service struct {
	Repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{Repo: repo}
}

func (s *Service) Record(input RecordInput) error {
	logEntry := &AuditLog{
		ActorUserID: input.ActorUserID,
		Action:      input.Action,
		Resource:    input.Resource,
		ResourceID:  input.ResourceID,
		Status:      input.Status,
		Description: input.Description,
		IPAddress:   input.IPAddress,
		UserAgent:   input.UserAgent,
	}

	return s.Repo.Create(logEntry)
}

func (s *Service) SafeRecord(input RecordInput) {
	if s == nil || s.Repo == nil {
		return
	}

	if err := s.Record(input); err != nil {
		log.Printf("audit log failed for %s/%s: %v", input.Resource, input.Action, err)
	}
}

func (s *Service) GetLogs(page, limit int, action, resource string) ([]AuditLog, int64, error) {
	return s.Repo.FindAllWithFilter(page, limit, action, resource)
}
