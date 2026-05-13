package tools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/bimross/skills-mcp-server/internal/skills"
)

type createSkillTool struct {
	store *skills.Store
}

func (t *createSkillTool) Name() string { return "create_skill" }

func (t *createSkillTool) Definition() map[string]any {
	return map[string]any{
		"name": "create_skill",
		"description": "Create a new Agent Skill on disk under SKILLS_MCP_SERVER_DIR: writes <name>/SKILL.md. " +
			"Only call when the user explicitly asked to add a new packaged skill with a clear kebab-case id (e.g. talk-better). " +
			"Requires name, description, and instructions; optional license, compatibility, metadata, allowedTools. " +
			"Destructive if misused—do not call for casual chat.",
		"inputSchema": map[string]any{
			"type":     "object",
			"required": []string{"name", "description", "instructions"},
			"properties": map[string]any{
				"name":          map[string]any{"type": "string", "description": "Skill id / directory name (lowercase kebab-case)."},
				"description":   map[string]any{"type": "string", "description": "Short catalog description (1–1024 chars)."},
				"instructions":  map[string]any{"type": "string", "description": "Main markdown body after frontmatter."},
				"license":       map[string]any{"type": "string"},
				"compatibility": map[string]any{"type": "string"},
				"metadata":      map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}},
				"allowedTools":  map[string]any{"type": "string", "description": "Optional; YAML frontmatter may use allowed-tools as alias."},
			},
		},
	}
}

func (t *createSkillTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(arguments, &raw); err != nil {
		return nil, err
	}
	input, err := decodeCreateOrUpdateInput(raw)
	if err != nil {
		return nil, err
	}
	return t.store.CreateSkill(strings.TrimSpace(input.Name), input)
}

type updateSkillTool struct {
	store *skills.Store
}

func (t *updateSkillTool) Name() string { return "update_skill" }

func (t *updateSkillTool) Definition() map[string]any {
	return map[string]any{
		"name": "update_skill",
		"description": "Replace an existing skill's SKILL.md from structured fields. " +
			"Only when the user clearly asked to edit that skill by id. name must match the existing skill folder. " +
			"Same fields as create_skill.",
		"inputSchema": map[string]any{
			"type":     "object",
			"required": []string{"name", "description", "instructions"},
			"properties": map[string]any{
				"name":          map[string]any{"type": "string"},
				"description":   map[string]any{"type": "string"},
				"instructions":  map[string]any{"type": "string"},
				"license":       map[string]any{"type": "string"},
				"compatibility": map[string]any{"type": "string"},
				"metadata":      map[string]any{"type": "object", "additionalProperties": map[string]any{"type": "string"}},
				"allowedTools":  map[string]any{"type": "string", "description": "Optional; YAML may use allowed-tools as alias."},
			},
		},
	}
}

func (t *updateSkillTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(arguments, &raw); err != nil {
		return nil, err
	}
	input, err := decodeCreateOrUpdateInput(raw)
	if err != nil {
		return nil, err
	}
	return t.store.UpdateSkill(strings.TrimSpace(input.Name), input)
}

type deleteSkillTool struct {
	store *skills.Store
}

func (t *deleteSkillTool) Name() string { return "delete_skill" }

func (t *deleteSkillTool) Definition() map[string]any {
	return map[string]any{
		"name": "delete_skill",
		"description": "Permanently remove a skill directory and all files under it. " +
			"Only call when the user explicitly confirmed they want to delete that skill id. " +
			"Requires name (skill id). Irreversible.",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"name"},
			"properties": map[string]any{"name": map[string]any{"type": "string", "description": "Skill id to delete."}},
		},
	}
}

func (t *deleteSkillTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(args.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	if err := t.store.DeleteSkill(name); err != nil {
		return nil, err
	}
	return map[string]any{"deleted": true, "name": name}, nil
}

// decodeCreateOrUpdateInput builds CreateOrUpdateSkillInput from JSON object keys (supports allowed-tools alias).
func decodeCreateOrUpdateInput(raw map[string]json.RawMessage) (skills.CreateOrUpdateSkillInput, error) {
	var input skills.CreateOrUpdateSkillInput
	b, err := json.Marshal(raw)
	if err != nil {
		return input, err
	}
	if err := json.Unmarshal(b, &input); err != nil {
		return input, err
	}
	if v, ok := raw["allowed-tools"]; ok && input.AllowedTools == "" {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			input.AllowedTools = s
		}
	}
	return input, nil
}
