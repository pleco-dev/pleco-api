package audit

import (
	"fmt"
	"strconv"
	"time"

	"go-api-starterkit/internal/httpx"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	AuditService *Service
}

func NewHandler(auditService *Service) *Handler {
	return &Handler{AuditService: auditService}
}

func (h *Handler) GetLogs(c *gin.Context) {
	filter, err := buildFilter(c)
	if err != nil {
		httpx.Error(c, 400, err.Error())
		return
	}

	logs, total, err := h.AuditService.GetLogs(filter)
	if err != nil {
		httpx.Error(c, 500, "Failed to fetch audit logs")
		return
	}

	httpx.Success(c, 200, "Audit logs fetched", logs, gin.H{
		"page":          filter.Page,
		"limit":         filter.Limit,
		"total":         total,
		"action":        filter.Action,
		"resource":      filter.Resource,
		"status":        filter.Status,
		"actor_user_id": filter.ActorUserID,
		"search":        filter.Search,
		"date_from":     formatTime(filter.DateFrom),
		"date_to":       formatTime(filter.DateTo),
	})
}

func (h *Handler) ExportLogs(c *gin.Context) {
	filter, err := buildFilter(c)
	if err != nil {
		httpx.Error(c, 400, err.Error())
		return
	}

	payload, err := h.AuditService.ExportLogsCSV(filter)
	if err != nil {
		httpx.Error(c, 500, "Failed to export audit logs")
		return
	}

	filename := fmt.Sprintf("audit-logs-%s.csv", time.Now().UTC().Format("20060102-150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filename))
	c.Data(200, "text/csv; charset=utf-8", payload)
}

func buildFilter(c *gin.Context) (Filter, error) {
	page := 1
	limit := 10

	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	filter := Filter{
		Page:     page,
		Limit:    limit,
		Action:   c.Query("action"),
		Resource: c.Query("resource"),
		Status:   c.Query("status"),
		Search:   c.Query("search"),
	}

	if actorID := c.Query("actor_user_id"); actorID != "" {
		parsed, err := strconv.ParseUint(actorID, 10, 64)
		if err != nil || parsed == 0 {
			return Filter{}, fmt.Errorf("actor_user_id must be a positive integer")
		}
		value := uint(parsed)
		filter.ActorUserID = &value
	}

	var err error
	if dateFrom := c.Query("date_from"); dateFrom != "" {
		parsed := time.Time{}
		parsed, err = time.Parse(time.RFC3339, dateFrom)
		if err != nil {
			return Filter{}, fmt.Errorf("date_from must use RFC3339 format")
		}
		filter.DateFrom = &parsed
	}
	if dateTo := c.Query("date_to"); dateTo != "" {
		parsed := time.Time{}
		parsed, err = time.Parse(time.RFC3339, dateTo)
		if err != nil {
			return Filter{}, fmt.Errorf("date_to must use RFC3339 format")
		}
		filter.DateTo = &parsed
	}
	if filter.DateFrom != nil && filter.DateTo != nil && filter.DateFrom.After(*filter.DateTo) {
		return Filter{}, fmt.Errorf("date_from must be before or equal to date_to")
	}

	return filter, nil
}

func formatTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return value.UTC().Format(time.RFC3339)
}
