package audit

type InvestigateRequest struct {
	Action      string `json:"action"`
	Resource    string `json:"resource"`
	Status      string `json:"status"`
	ActorUserID *uint  `json:"actor_user_id"`
	Search      string `json:"search"`
	DateFrom    string `json:"date_from"`
	DateTo      string `json:"date_to"`
	Limit       int    `json:"limit"`
}

type InvestigationResult struct {
	Summary           string   `json:"summary"`
	Timeline          []string `json:"timeline"`
	SuspiciousSignals []string `json:"suspicious_signals"`
	Recommendations   []string `json:"recommendations"`
}
