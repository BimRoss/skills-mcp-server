package tools

import (
	"context"
	"encoding/json"

	"github.com/bimross/skills-mcp-server/internal/readweb"
)

type readWebTool struct {
	client *readweb.Client
}

func (t *readWebTool) Name() string { return "read_web" }

func (t *readWebTool) Definition() map[string]any {
	return map[string]any{
		"name":        "read_web",
		"description": "Run internet research and return summary with citations",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"query"},
			"properties": map[string]any{"query": map[string]any{"type": "string"}},
		},
	}
}

func (t *readWebTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	result, err := t.client.Run(ctx, args.Query)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"fallbackText": result.Summary,
		"finalSummary": result.Summary,
		"citations":    result.Citations,
	}, nil
}
