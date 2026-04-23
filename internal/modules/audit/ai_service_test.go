package audit

import (
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
)

func TestParseInvestigationResultSupportsStructuredValues(t *testing.T) {
	raw := `{
	  "summary": {
	    "total_actions": 2,
	    "common_resource": "auth"
	  },
	  "timeline": [
	    {
	      "time": "2026-04-22T09:10:00Z",
	      "action": "login"
	    }
	  ],
	  "suspicious_signals": [
	    {
	      "pattern": "invalid credentials",
	      "total_occurrences": 2
	    }
	  ],
	  "recommendations": [
	    {
	      "action": "Review account activity"
	    }
	  ]
	}`

	result, err := parseInvestigationResult(raw)
	if err != nil {
		t.Fatalf("expected parser to succeed, got error: %v", err)
	}

	if result.Summary == "" {
		t.Fatal("expected summary to be stringified")
	}
	if !strings.Contains(result.Timeline[0], "time: 2026-04-22T09:10:00Z") {
		t.Fatalf("expected human-readable timeline item, got %q", result.Timeline[0])
	}
	if !strings.Contains(result.SuspiciousSignals[0], "pattern: invalid credentials") {
		t.Fatalf("expected human-readable suspicious signal, got %q", result.SuspiciousSignals[0])
	}
	if !strings.Contains(result.Recommendations[0], "action: Review account activity") {
		t.Fatalf("expected human-readable recommendation, got %q", result.Recommendations[0])
	}
	if len(result.Timeline) != 1 {
		t.Fatalf("expected 1 timeline item, got %d", len(result.Timeline))
	}
	if len(result.SuspiciousSignals) != 1 {
		t.Fatalf("expected 1 suspicious signal, got %d", len(result.SuspiciousSignals))
	}
	if len(result.Recommendations) != 1 {
		t.Fatalf("expected 1 recommendation, got %d", len(result.Recommendations))
	}
}

func TestBuildInvestigationContextSeparatesActorAndIPAddress(t *testing.T) {
	logs := []AuditLog{
		{
			Model: gorm.Model{
				ID:        7,
				CreatedAt: time.Date(2026, 4, 22, 2, 31, 14, 0, time.UTC),
			},
			Action:      "login",
			Resource:    "auth",
			Status:      "failed",
			Description: "invalid credentials",
			IPAddress:   "::1",
		},
	}

	context := buildInvestigationContext(logs)

	if !strings.Contains(context, "actor_user_id: n/a") {
		t.Fatalf("expected missing actor to be explicit, got %q", context)
	}
	if !strings.Contains(context, "ip_address: ::1") {
		t.Fatalf("expected ip_address field to be explicit, got %q", context)
	}
	if !strings.Contains(context, "description: \"invalid credentials\"") {
		t.Fatalf("expected description to be included, got %q", context)
	}
	if !strings.Contains(context, "Overview: failed_logs=1 unique_ip_addresses=1 unique_actor_user_ids=0") {
		t.Fatalf("expected overview summary in context, got %q", context)
	}
}

func TestMergeInvestigationHeuristicsBuildsReadableFallbacks(t *testing.T) {
	logs := []AuditLog{
		{
			Model: gorm.Model{
				ID:        1,
				CreatedAt: time.Date(2026, 4, 22, 2, 31, 14, 0, time.UTC),
			},
			Action:    "login",
			Resource:  "auth",
			Status:    "failed",
			IPAddress: "::1",
		},
		{
			Model: gorm.Model{
				ID:        2,
				CreatedAt: time.Date(2026, 4, 22, 2, 32, 14, 0, time.UTC),
			},
			Action:    "login",
			Resource:  "auth",
			Status:    "failed",
			IPAddress: "::1",
		},
	}

	result := mergeInvestigationHeuristics(&InvestigationResult{
		Timeline:          []string{"2026-04-22T02:31:14Z"},
		SuspiciousSignals: []string{"None detected."},
	}, logs)

	if len(result.Timeline) == 0 || !strings.Contains(result.Timeline[0], "failed login on auth from ip ::1") {
		t.Fatalf("expected readable timeline fallback, got %#v", result.Timeline)
	}
	if len(result.SuspiciousSignals) == 0 || !strings.Contains(result.SuspiciousSignals[0], "failed events originated from ip ::1") {
		t.Fatalf("expected heuristic suspicious signal, got %#v", result.SuspiciousSignals)
	}
	if len(result.Recommendations) == 0 {
		t.Fatal("expected derived recommendations")
	}
}

func TestParseInvestigationResultExtractsJSONFromPreamble(t *testing.T) {
	raw := `Here is the investigation result:
{
  "summary": "Failed login attempt detected.",
  "timeline": ["2026-04-22T02:31:14Z - failed login on auth from ip ::1 (actor_user_id: n/a)"],
  "suspicious_signals": ["Repeated failed logins from ip ::1."],
  "recommendations": ["Review the affected account."]
}`

	result, err := parseInvestigationResult(raw)
	if err != nil {
		t.Fatalf("expected parser to extract JSON, got error: %v", err)
	}
	if result.Summary != "Failed login attempt detected." {
		t.Fatalf("unexpected summary: %q", result.Summary)
	}
}

func TestParseInvestigationResultExtractsJSONFromCodeFence(t *testing.T) {
	raw := "```json\n{\"summary\":\"ok\",\"timeline\":[\"t\"],\"suspicious_signals\":[\"s\"],\"recommendations\":[\"r\"]}\n```"

	result, err := parseInvestigationResult(raw)
	if err != nil {
		t.Fatalf("expected parser to handle code fence, got error: %v", err)
	}
	if result.Summary != "ok" {
		t.Fatalf("unexpected summary: %q", result.Summary)
	}
}
