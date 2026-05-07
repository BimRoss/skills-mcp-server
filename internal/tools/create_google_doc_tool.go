package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/bimross/skills-mcp-server/internal/googledocs"
)

type createGoogleDocTool struct {
	cfg googledocs.EnvConfig
}

func (t *createGoogleDocTool) Name() string { return "create_google_doc" }

func (t *createGoogleDocTool) Definition() map[string]any {
	return map[string]any{
		"name":        "create_google_doc",
		"description": "Create a Google Doc, insert body text, and grant Drive sharing to listed Google accounts. Use when the user wants a shared Google document, meeting notes captured in Docs, or a draft others can edit. Requires server OAuth refresh token + Docs/Drive scopes. Editors/commenters/viewers must be valid email addresses.",
		"inputSchema": map[string]any{
			"type":     "object",
			"required": []string{"intent", "title", "editors"},
			"properties": map[string]any{
				"intent": map[string]any{
					"type":        "string",
					"description": "Short natural-language description of why the doc exists; used as document body when body is omitted.",
				},
				"title": map[string]any{
					"type":        "string",
					"description": "Google Doc title.",
				},
				"editors": map[string]any{
					"type":        "array",
					"description": "Email addresses to grant writer (editor) access.",
					"items":       map[string]any{"type": "string"},
				},
				"commenters": map[string]any{
					"type":        "array",
					"description": "Optional emails for commenter access (not also editors).",
					"items":       map[string]any{"type": "string"},
				},
				"viewers": map[string]any{
					"type":        "array",
					"description": "Optional read-only viewer emails (excluding editors and commenters).",
					"items":       map[string]any{"type": "string"},
				},
				"type": map[string]any{
					"type":        "string",
					"description": "Optional document style hint for callers (ignored by this server implementation).",
				},
				"length": map[string]any{
					"type":        "string",
					"description": "Optional length hint for callers (ignored by this server implementation).",
				},
				"body": map[string]any{
					"type":        "string",
					"description": "Plain text document body. When empty, intent is used as the body.",
				},
			},
		},
	}
}

func (t *createGoogleDocTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Intent     string   `json:"intent"`
		Title      string   `json:"title"`
		Editors    []string `json:"editors"`
		Commenters []string `json:"commenters"`
		Viewers    []string `json:"viewers"`
		Body       string   `json:"body"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	cfg := t.cfg
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(args.Title) == "" {
		return nil, fmt.Errorf("create_google_doc: title is required")
	}
	if len(args.Editors) == 0 {
		return nil, fmt.Errorf("create_google_doc: editors must include at least one email")
	}
	body := strings.TrimSpace(args.Body)
	if body == "" {
		body = strings.TrimSpace(args.Intent)
	}
	if body == "" {
		return nil, fmt.Errorf("create_google_doc: body is empty and intent is empty")
	}

	client, err := googledocs.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	createRes, err := client.Create(ctx, googledocs.CreateInput{
		Title: strings.TrimSpace(args.Title),
		Body:  body,
	})
	if err != nil {
		return nil, err
	}
	editors := googledocs.DedupeEmails(args.Editors)
	for _, email := range editors {
		if err := client.GrantEditor(ctx, createRes.DocumentID, email); err != nil {
			return nil, err
		}
	}
	for _, email := range googledocs.SubtractEmails(googledocs.DedupeEmails(args.Commenters), args.Editors) {
		if err := client.GrantCommenter(ctx, createRes.DocumentID, email); err != nil {
			return nil, err
		}
	}
	viewers := googledocs.SubtractEmails(googledocs.DedupeEmails(args.Viewers), args.Editors)
	viewers = googledocs.SubtractEmails(viewers, args.Commenters)
	for _, email := range viewers {
		if err := client.GrantViewer(ctx, createRes.DocumentID, email); err != nil {
			return nil, err
		}
	}

	summary := "create_google_doc completed"
	fallback := "Created Google Doc: " + strings.TrimSpace(createRes.URL)
	return map[string]any{
		"fallbackText": fallback,
		"finalSummary": summary,
		"docId":        createRes.DocumentID,
		"docUrl":       createRes.URL,
	}, nil
}

// RunCreateGoogleDoc runs the same logic as the create_google_doc MCP tool (for REST aliases).
func RunCreateGoogleDoc(ctx context.Context, cfg googledocs.EnvConfig, arguments json.RawMessage) (any, error) {
	return (&createGoogleDocTool{cfg: cfg}).Call(ctx, arguments)
}
