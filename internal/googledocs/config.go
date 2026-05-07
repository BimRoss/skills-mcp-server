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

// LoadFromEnv reads Google OAuth env vars for create_google_doc.
// Joanne keys win first so a combined .env (portal GOOGLE_OAUTH_* + worker JOANNE_*)
// does not pair a portal client with a Joanne-issued refresh token (oauth2 unauthorized_client).
func LoadFromEnv() EnvConfig {
	return EnvConfig{
		ClientID: firstNonEmpty(
			os.Getenv("JOANNE_GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_CLIENT_ID"),
			os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		),
		ClientSecret: firstNonEmpty(
			os.Getenv("JOANNE_GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_CLIENT_SECRET"),
			os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		),
		RefreshToken: firstNonEmpty(
			os.Getenv("JOANNE_GOOGLE_REFRESH_TOKEN"),
			os.Getenv("GOOGLE_REFRESH_TOKEN"),
			os.Getenv("GOOGLE_OAUTH_REFRESH_TOKEN"),
		),
	}
}

// Validate returns an error if any required field is missing.
func (c EnvConfig) Validate() error {
	if strings.TrimSpace(c.ClientID) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth client id (set JOANNE_GOOGLE_CLIENT_ID, or GOOGLE_CLIENT_ID / GOOGLE_OAUTH_CLIENT_ID)")
	}
	if strings.TrimSpace(c.ClientSecret) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth client secret (set JOANNE_GOOGLE_CLIENT_SECRET, or GOOGLE_CLIENT_SECRET / GOOGLE_OAUTH_CLIENT_SECRET)")
	}
	if strings.TrimSpace(c.RefreshToken) == "" {
		return fmt.Errorf("create_google_doc: missing Google OAuth refresh token (set JOANNE_GOOGLE_REFRESH_TOKEN, or GOOGLE_REFRESH_TOKEN / GOOGLE_OAUTH_REFRESH_TOKEN; Docs + drive.file)")
	}
	return nil
}
