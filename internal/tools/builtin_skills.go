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
		"description": "Return packaged Agent Skill entries (folder id, description, metadata) from disk—this is the SKILL.md catalog, not the list of MCP executable tools. Use for skimmable catalogs, browsing packaged guidance, or before read_skill when you need ids. Optional `query` filters by substring over name/description. When the user asks what MCP tools or invocable tools exist, summarize from the tools you were given in this session; call this when they want the packaged skill library (or both, if unclear).",
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
		"description": "Find skills whose name or description match `query`. Use when the user describes a goal without naming a skill id (e.g. roster, directory, onboarding, research), or when you need the canonical skill folder name before calling read_skill.",
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
		"description": "Load full markdown instructions for one packaged skill. Required argument `name` is the skill id (folder name from list_skills or search_skills). Use when the user names a skill, asks to follow packaged guidance, or when answers should come from stored SKILL.md text—not free recall alone. For directory/roster-style questions, use search_skills with relevant keywords to find the deployment's skill id first; ids are deployment-specific—never present a skill folder name as if it were an MCP tool name.",
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
