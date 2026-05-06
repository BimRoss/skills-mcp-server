package app

import (
	"fmt"
	"net/http"

	"github.com/bimross/skills-mcp-server/internal/config"
	"github.com/bimross/skills-mcp-server/internal/httpapi"
	"github.com/bimross/skills-mcp-server/internal/mcp"
	"github.com/bimross/skills-mcp-server/internal/readweb"
	"github.com/bimross/skills-mcp-server/internal/skills"
)

type App struct {
	handler http.Handler
}

func New(cfg config.Config) (*App, error) {
	store, err := skills.NewStore(cfg.SkillsDir)
	if err != nil {
		return nil, fmt.Errorf("build store: %w", err)
	}
	if cfg.SeedExamples {
		if err := skills.SeedExamplesFromDir(store, cfg.ExamplesDir); err != nil {
			return nil, fmt.Errorf("seed examples into skills dir: %w", err)
		}
	}
	readWeb := readweb.New(readweb.Config{
		APIKey:            cfg.GeminiAPIKey,
		Model:             cfg.GeminiModel,
		EnableWebResearch: cfg.EnableWebResearch,
	})
	mux := http.NewServeMux()
	httpapi.New(store, readWeb).Register(mux)
	mcp.New(store, readWeb).Register(mux)
	return &App{handler: mux}, nil
}

func (a *App) Handler() http.Handler {
	return a.handler
}
