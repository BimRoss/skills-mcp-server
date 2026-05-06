package readweb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Config struct {
	APIKey            string
	Model             string
	EnableWebResearch bool
}

type Result struct {
	Query     string   `json:"query"`
	Summary   string   `json:"summary"`
	Citations []string `json:"citations"`
}

type Client struct {
	cfg Config
}

func New(cfg Config) *Client {
	return &Client{cfg: cfg}
}

func (c *Client) Run(ctx context.Context, query string) (Result, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return Result{}, fmt.Errorf("query is required")
	}
	if strings.TrimSpace(c.cfg.APIKey) == "" {
		return Result{}, fmt.Errorf("missing GEMINI_API_KEY")
	}
	if strings.TrimSpace(c.cfg.Model) == "" {
		return Result{}, fmt.Errorf("missing GEMINI_MODEL")
	}

	requestBody := map[string]any{
		"contents": []any{
			map[string]any{
				"parts": []any{
					map[string]any{
						"text": "Research the following query and respond with a concise summary. Do not include source links in the text body.\n\nQuery: " + query,
					},
				},
			},
		},
		"generationConfig": map[string]any{
			"temperature":      0.2,
			"responseMimeType": "text/plain",
		},
	}
	if c.cfg.EnableWebResearch {
		requestBody["tools"] = []any{
			map[string]any{"google_search": map[string]any{}},
		}
	}

	payload, err := json.Marshal(requestBody)
	if err != nil {
		return Result{}, fmt.Errorf("marshal request: %w", err)
	}
	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		url.PathEscape(c.cfg.Model),
		url.QueryEscape(c.cfg.APIKey),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return Result{}, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 45 * time.Second}).Do(req)
	if err != nil {
		return Result{}, fmt.Errorf("gemini request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{}, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return Result{}, fmt.Errorf("gemini returned status %d", resp.StatusCode)
	}

	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
			GroundingMetadata struct {
				GroundingChunks []struct {
					Web struct {
						URI string `json:"uri"`
					} `json:"web"`
				} `json:"groundingChunks"`
			} `json:"groundingMetadata"`
		} `json:"candidates"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return Result{}, fmt.Errorf("parse response: %w", err)
	}
	summary := ""
	if len(parsed.Candidates) > 0 && len(parsed.Candidates[0].Content.Parts) > 0 {
		summary = strings.TrimSpace(parsed.Candidates[0].Content.Parts[0].Text)
	}
	if summary == "" {
		return Result{}, fmt.Errorf("empty gemini summary")
	}

	citationSet := map[string]struct{}{}
	citations := make([]string, 0)
	for _, cand := range parsed.Candidates {
		for _, chunk := range cand.GroundingMetadata.GroundingChunks {
			link := strings.TrimSpace(chunk.Web.URI)
			if link == "" {
				continue
			}
			if _, exists := citationSet[link]; exists {
				continue
			}
			citationSet[link] = struct{}{}
			citations = append(citations, link)
		}
	}
	return Result{
		Query:     query,
		Summary:   summary,
		Citations: citations,
	}, nil
}
