package tools

import (
	"context"
	"encoding/base64"
	"encoding/json"

	"github.com/bimross/skills-mcp-server/internal/skills"
)

type listSkillsTool struct {
	store *skills.Store
}

func (t *listSkillsTool) Name() string { return "list_skills" }

func (t *listSkillsTool) Definition() map[string]any {
	return map[string]any{
		"name":        "list_skills",
		"description": "List all available skills",
		"inputSchema": map[string]any{
			"type":       "object",
			"properties": map[string]any{"query": map[string]any{"type": "string"}},
		},
	}
}

func (t *listSkillsTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Query string `json:"query"`
	}
	_ = json.Unmarshal(arguments, &args)
	items, err := t.store.ListSkills(args.Query)
	if err != nil {
		return nil, err
	}
	return items, nil
}

type searchSkillsTool struct {
	store *skills.Store
}

func (t *searchSkillsTool) Name() string { return "search_skills" }

func (t *searchSkillsTool) Definition() map[string]any {
	return map[string]any{
		"name":        "search_skills",
		"description": "Search skills by text query",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"query"},
			"properties": map[string]any{"query": map[string]any{"type": "string"}},
		},
	}
}

func (t *searchSkillsTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	items, err := t.store.ListSkills(args.Query)
	if err != nil {
		return nil, err
	}
	return items, nil
}

type readSkillTool struct {
	store *skills.Store
}

func (t *readSkillTool) Name() string { return "read_skill" }

func (t *readSkillTool) Definition() map[string]any {
	return map[string]any{
		"name":        "read_skill",
		"description": "Read full content for a single skill",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"name"},
			"properties": map[string]any{"name": map[string]any{"type": "string"}},
		},
	}
}

func (t *readSkillTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	return t.store.ReadSkill(args.Name)
}

type listSkillResourcesTool struct {
	store *skills.Store
}

func (t *listSkillResourcesTool) Name() string { return "list_skill_resources" }

func (t *listSkillResourcesTool) Definition() map[string]any {
	return map[string]any{
		"name":        "list_skill_resources",
		"description": "List resources under scripts/references/assets for a skill",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"name"},
			"properties": map[string]any{"name": map[string]any{"type": "string"}},
		},
	}
}

func (t *listSkillResourcesTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	return t.store.ListResources(args.Name)
}

type readSkillResourceTool struct {
	store *skills.Store
}

func (t *readSkillResourceTool) Name() string { return "read_skill_resource" }

func (t *readSkillResourceTool) Definition() map[string]any {
	return map[string]any{
		"name":        "read_skill_resource",
		"description": "Read a resource file content by path",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"name", "path"},
			"properties": map[string]any{"name": map[string]any{"type": "string"}, "path": map[string]any{"type": "string"}},
		},
	}
}

func (t *readSkillResourceTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	content, info, err := t.store.ReadResource(args.Name, args.Path)
	if err != nil {
		return nil, err
	}
	isText := skills.IsLikelyText(info.Path)
	encoded := ""
	if isText {
		encoded = string(content)
	} else {
		encoded = base64.StdEncoding.EncodeToString(content)
	}
	return map[string]any{
		"path":      info.Path,
		"sizeBytes": info.SizeBytes,
		"updatedAt": info.UpdatedAt,
		"encoding":  encodingLabel(isText),
		"content":   encoded,
	}, nil
}

type getSkillResourceInfoTool struct {
	store *skills.Store
}

func (t *getSkillResourceInfoTool) Name() string { return "get_skill_resource_info" }

func (t *getSkillResourceInfoTool) Definition() map[string]any {
	return map[string]any{
		"name":        "get_skill_resource_info",
		"description": "Get metadata for a resource without content",
		"inputSchema": map[string]any{
			"type":       "object",
			"required":   []string{"name", "path"},
			"properties": map[string]any{"name": map[string]any{"type": "string"}, "path": map[string]any{"type": "string"}},
		},
	}
}

func (t *getSkillResourceInfoTool) Call(ctx context.Context, arguments json.RawMessage) (any, error) {
	var args struct {
		Name string `json:"name"`
		Path string `json:"path"`
	}
	if err := json.Unmarshal(arguments, &args); err != nil {
		return nil, err
	}
	_, info, err := t.store.ReadResource(args.Name, args.Path)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func encodingLabel(isText bool) string {
	if isText {
		return "text"
	}
	return "base64"
}
