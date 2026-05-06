package tools

import (
	"encoding/json"
)

// MCPResult wraps structured tool output for MCP tools/call (text mirror + structuredContent).
func MCPResult(value any) map[string]any {
	b, _ := json.Marshal(value)
	return map[string]any{
		"content": []map[string]any{
			{
				"type": "text",
				"text": string(b),
			},
		},
		"structuredContent": value,
	}
}
