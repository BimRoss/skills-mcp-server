package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port                    string
	SkillsDir               string
	GeminiAPIKey            string
	GeminiModel             string
	EnableWebResearch       bool
	SeedBuiltinReadWebSkill bool
}

func FromEnv() Config {
	return Config{
		Port:                    envOrDefault("SKILLS_MCP_SERVER_PORT", "8081"),
		SkillsDir:               envOrDefault("SKILLS_MCP_SERVER_DIR", "./skills"),
		GeminiAPIKey:            os.Getenv("GEMINI_API_KEY"),
		GeminiModel:             envOrDefault("GEMINI_MODEL", "gemini-2.5-flash"),
		EnableWebResearch:       strings.EqualFold(envOrDefault("ENABLE_WEB_RESEARCH", "true"), "true"),
		SeedBuiltinReadWebSkill: strings.EqualFold(envOrDefault("SKILLS_SEED_BUILTIN_READ_WEB", "true"), "true"),
	}
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
