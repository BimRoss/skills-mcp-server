package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Port              string
	SkillsDir         string
	GeminiAPIKey      string
	GeminiModel       string
	EnableWebResearch bool
	SeedExamples      bool
	ExamplesDir       string
}

func FromEnv() Config {
	return Config{
		Port:              envOrDefault("SKILLS_MCP_SERVER_PORT", "8081"),
		SkillsDir:         envOrDefault("SKILLS_MCP_SERVER_DIR", "./skills"),
		GeminiAPIKey:      os.Getenv("GEMINI_API_KEY"),
		GeminiModel:       envOrDefault("GEMINI_MODEL", "gemini-2.5-flash"),
		EnableWebResearch: strings.EqualFold(envOrDefault("ENABLE_WEB_RESEARCH", "true"), "true"),
		SeedExamples:      seedExamplesDefaultTrue(),
		ExamplesDir:       examplesDirFromEnv(),
	}
}

func seedExamplesDefaultTrue() bool {
	v := strings.TrimSpace(os.Getenv("SKILLS_SEED_EXAMPLES"))
	if v == "" {
		return true
	}
	return strings.EqualFold(v, "true") || v == "1"
}

// examplesDirFromEnv: explicit SKILLS_EXAMPLES_DIR, else ./examples for local dev (empty disables only if SeedExamples false).
func examplesDirFromEnv() string {
	if v := strings.TrimSpace(os.Getenv("SKILLS_EXAMPLES_DIR")); v != "" {
		return v
	}
	return filepath.Clean("examples")
}

func (c Config) ListenAddress() string {
	return fmt.Sprintf(":%s", c.Port)
}

func envOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
