package googledocs

import (
	"fmt"
	"os"
	"strings"
)

// EnvConfig holds OAuth refresh-token credentials for Google Docs + Drive file scope.
type EnvConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

// LoadFromEnv reads global Google OAuth env vars. Resolution order prefers explicit
// GOOGLE_CLIENT_* (agent-factory / employee-factory), then GOOGLE_OAUTH_* (portal-style),
// then JOANNE_GOOGLE_* employee overrides.
func LoadFromEnv() EnvConfig {
	return EnvConfig{
		ClientID: firstNonEmpty(
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
			os.Getenv("JOANNE_GOOGLE_CLIENT_ID"),
		),
		ClientSecret: firstNonEmpty(
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
			os.Getenv("JOANNE_GOOGLE_CLIENT_SECRET"),
		),
		RefreshToken: firstNonEmpty(
			os.Getenv("GOOGLE_REFRESH_TOKEN"),
			os.Getenv("GOOGLE_OAUTH_REFRESH_TOKEN"),
			os.Getenv("JOANNE_GOOGLE_REFRESH_TOKEN"),
		),
	}
}

// Validate returns an error if any required field is missing.
func (c EnvConfig) Validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth client id (set GOOGLE_CLIENT_ID or GOOGLE_OAUTH_CLIENT_ID)")
	}
	if strings.TrimSpace(c.ClientSecret) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth client secret (set GOOGLE_CLIENT_SECRET or GOOGLE_OAUTH_CLIENT_SECRET)")
	}
	if strings.TrimSpace(c.RefreshToken) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth refresh token (set GOOGLE_REFRESH_TOKEN or GOOGLE_OAUTH_REFRESH_TOKEN; user-delegated Docs scope)")
	}
	return nil
}
