package tests

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-api-starterkit/internal/modules/audit"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type stubAuditRepo struct {
	findAllWithFilter func(filter audit.Filter) ([]audit.AuditLog, int64, error)
	findForExport     func(filter audit.Filter) ([]audit.AuditLog, error)
}

func (s *stubAuditRepo) Create(log *audit.AuditLog) error {
	return nil
}

func (s *stubAuditRepo) FindAllWithFilter(filter audit.Filter) ([]audit.AuditLog, int64, error) {
	if s.findAllWithFilter != nil {
		return s.findAllWithFilter(filter)
	}
	return nil, 0, nil
}

func (s *stubAuditRepo) FindForExport(filter audit.Filter) ([]audit.AuditLog, error) {
	if s.findForExport != nil {
		return s.findForExport(filter)
	}
	return nil, nil
}

func TestAuditHandler_GetLogs_WithExtendedFilter(t *testing.T) {
	expectedActorID := uint(42)
	expectedFrom, _ := time.Parse(time.RFC3339, "2026-04-20T00:00:00Z")
	expectedTo, _ := time.Parse(time.RFC3339, "2026-04-21T00:00:00Z")

	service := audit.NewService(&stubAuditRepo{
		findAllWithFilter: func(filter audit.Filter) ([]audit.AuditLog, int64, error) {
			assert.Equal(t, 2, filter.Page)
			assert.Equal(t, 20, filter.Limit)
			assert.Equal(t, "login", filter.Action)
			assert.Equal(t, "auth", filter.Resource)
			assert.Equal(t, "failed", filter.Status)
			assert.Equal(t, "203.0.113.1", filter.Search)
			assert.NotNil(t, filter.ActorUserID)
			assert.Equal(t, expectedActorID, *filter.ActorUserID)
			assert.NotNil(t, filter.DateFrom)
			assert.NotNil(t, filter.DateTo)
			assert.True(t, filter.DateFrom.Equal(expectedFrom))
			assert.True(t, filter.DateTo.Equal(expectedTo))

			return []audit.AuditLog{{Action: "login", Resource: "auth", Status: "failed"}}, 1, nil
		},
	})

	handler := audit.NewHandler(service)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/audit-logs?page=2&limit=20&action=login&resource=auth&status=failed&actor_user_id=42&search=203.0.113.1&date_from=2026-04-20T00:00:00Z&date_to=2026-04-21T00:00:00Z",
		nil,
	)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	bodyMap := decodeBodyMap(t, w)
	assert.Equal(t, "success", bodyMap["status"])
	assert.Equal(t, "Audit logs fetched", bodyMap["message"])
	meta := bodyMap["meta"].(map[string]interface{})
	assert.Equal(t, float64(2), meta["page"])
	assert.Equal(t, float64(20), meta["limit"])
	assert.Equal(t, "failed", meta["status"])
}

func TestAuditHandler_GetLogs_InvalidDateRange(t *testing.T) {
	handler := audit.NewHandler(audit.NewService(&stubAuditRepo{}))

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(
		http.MethodGet,
		"/audit-logs?date_from=2026-04-21T00:00:00Z&date_to=2026-04-20T00:00:00Z",
		nil,
	)

	handler.GetLogs(c)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	bodyMap := decodeBodyMap(t, w)
	assert.Equal(t, "error", bodyMap["status"])
	assert.Equal(t, "date_from must be before or equal to date_to", bodyMap["message"])
}

func TestAuditHandler_ExportLogs_ReturnsCSV(t *testing.T) {
	expectedActorID := uint(7)
	service := audit.NewService(&stubAuditRepo{
		findForExport: func(filter audit.Filter) ([]audit.AuditLog, error) {
			assert.NotNil(t, filter.ActorUserID)
			assert.Equal(t, expectedActorID, *filter.ActorUserID)
			return []audit.AuditLog{
				{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC),
					},
					ActorUserID: &expectedActorID,
					Action:      "login",
					Resource:    "auth",
					Status:      "success",
					Description: "user logged in",
					IPAddress:   "127.0.0.1",
					UserAgent:   "Postman",
				},
			}, nil
		},
	})
	handler := audit.NewHandler(service)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/audit-logs/export?actor_user_id=7", nil)

	handler.ExportLogs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Header().Get("Content-Type"), "text/csv")
	assert.Contains(t, w.Header().Get("Content-Disposition"), "attachment; filename=")
	assert.Contains(t, w.Body.String(), "id,created_at,actor_user_id,action,resource,resource_id,status,description,ip_address,user_agent")
	assert.Contains(t, w.Body.String(), "login,auth")
}

func TestAuditService_ExportLogsCSV_PropagatesRepositoryError(t *testing.T) {
	service := audit.NewService(&stubAuditRepo{
		findForExport: func(filter audit.Filter) ([]audit.AuditLog, error) {
			return nil, errors.New("db down")
		},
	})

	payload, err := service.ExportLogsCSV(audit.Filter{})

	assert.Nil(t, payload)
	assert.EqualError(t, err, "db down")
}
