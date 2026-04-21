package audit

import (
	"bytes"
	"encoding/csv"
	"log"
	"strconv"
)

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

func (s *Service) GetLogs(filter Filter) ([]AuditLog, int64, error) {
	return s.Repo.FindAllWithFilter(filter)
}

func (s *Service) ExportLogsCSV(filter Filter) ([]byte, error) {
	logs, err := s.Repo.FindForExport(filter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	if err := writer.Write([]string{
		"id",
		"created_at",
		"actor_user_id",
		"action",
		"resource",
		"resource_id",
		"status",
		"description",
		"ip_address",
		"user_agent",
	}); err != nil {
		return nil, err
	}

	for _, logEntry := range logs {
		if err := writer.Write([]string{
			strconv.FormatUint(uint64(logEntry.ID), 10),
			logEntry.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			uintPointerString(logEntry.ActorUserID),
			logEntry.Action,
			logEntry.Resource,
			uintPointerString(logEntry.ResourceID),
			logEntry.Status,
			logEntry.Description,
			logEntry.IPAddress,
			logEntry.UserAgent,
		}); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func uintPointerString(value *uint) string {
	if value == nil {
		return ""
	}
	return strconv.FormatUint(uint64(*value), 10)
}
