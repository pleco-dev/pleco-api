package config

import "strings"

type EmailConfig struct {
	APIKey      string
	From        string
	AppBaseURL  string
	FrontendURL string
}

type SocialConfig struct {
	GoogleClientID    string
	FacebookAppID     string
	FacebookAppSecret string
	AppleClientID     string
}

type AppConfig struct {
	Port              string
	DatabaseURL       string
	JWTSecret         []byte
	AdminEmail        string
	AdminPassword     string
	AutoRunMigrations bool
	AutoRunSeeds      bool
	Email             EmailConfig
	Social            SocialConfig
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port:              GetEnv("PORT", "8080"),
		DatabaseURL:       GetEnv("DATABASE_URL", ""),
		JWTSecret:         mustSecret("JWT_SECRET"),
		AdminEmail:        GetEnv("ADMIN_EMAIL", ""),
		AdminPassword:     GetEnv("ADMIN_PASSWORD", ""),
		AutoRunMigrations: envBool("AUTO_RUN_MIGRATIONS"),
		AutoRunSeeds:      envBool("AUTO_RUN_SEEDS"),
		Email: EmailConfig{
			APIKey:      GetEnv("SENDGRID_API_KEY", ""),
			From:        GetEnv("SENDGRID_EMAIL", ""),
			AppBaseURL:  firstNonEmptyEnv("APP_BASE_URL", "RENDER_EXTERNAL_URL", "http://localhost:8080"),
			FrontendURL: GetEnv("FRONTEND_URL", ""),
		},
		Social: SocialConfig{
			GoogleClientID:    GetEnv("GOOGLE_CLIENT_ID", ""),
			FacebookAppID:     GetEnv("FACEBOOK_APP_ID", ""),
			FacebookAppSecret: GetEnv("FACEBOOK_APP_SECRET", ""),
			AppleClientID:     GetEnv("APPLE_CLIENT_ID", ""),
		},
	}
}

func envBool(key string) bool {
	value := strings.TrimSpace(strings.ToLower(GetEnv(key, "")))
	return value == "1" || value == "true" || value == "yes"
}

func firstNonEmptyEnv(keys ...string) string {
	last := ""
	for _, key := range keys {
		if strings.Contains(key, "://") {
			last = key
			continue
		}
		if value := GetEnv(key, ""); value != "" {
			return value
		}
	}
	return last
}
