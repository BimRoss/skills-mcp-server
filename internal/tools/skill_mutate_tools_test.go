package tools

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/bimross/skills-mcp-server/internal/skills"
)

func TestCreateUpdateDeleteSkillTools(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	store, err := skills.NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	create := &createSkillTool{store: store}
	upd := &updateSkillTool{store: store}
	del := &deleteSkillTool{store: store}

	args := func(v any) json.RawMessage {
		b, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		return b
	}

	// invalid name
	_, err = create.Call(ctx, args(map[string]any{
		"name": "Bad_Name", "description": "d", "instructions": "i",
	}))
	if err == nil {
		t.Fatal("expected error for invalid name")
	}

	// create ok
	out, err := create.Call(ctx, args(map[string]any{
		"name": "alpha-skill", "description": "Alpha test skill", "instructions": "Do the thing.",
	}))
	if err != nil {
		t.Fatal(err)
	}
	skill, ok := out.(skills.Skill)
	if !ok || skill.Name != "alpha-skill" {
		t.Fatalf("unexpected create result: %#v", out)
	}

	// duplicate create
	_, err = create.Call(ctx, args(map[string]any{
		"name": "alpha-skill", "description": "x", "instructions": "y",
	}))
	if err == nil {
		t.Fatal("expected error on duplicate create")
	}

	// update ok
	out, err = upd.Call(ctx, args(map[string]any{
		"name": "alpha-skill", "description": "Alpha v2", "instructions": "Updated body.",
	}))
	if err != nil {
		t.Fatal(err)
	}
	skill = out.(skills.Skill)
	if skill.Description != "Alpha v2" || skill.Instructions != "Updated body." {
		t.Fatalf("update mismatch: %+v", skill)
	}

	// update missing
	_, err = upd.Call(ctx, args(map[string]any{
		"name": "nope-skill", "description": "d", "instructions": "i",
	}))
	if err == nil {
		t.Fatal("expected error updating missing skill")
	}

	// delete ok
	out, err = del.Call(ctx, args(map[string]any{"name": "alpha-skill"}))
	if err != nil {
		t.Fatal(err)
	}
	m := out.(map[string]any)
	if m["deleted"] != true || m["name"] != "alpha-skill" {
		t.Fatalf("delete result: %#v", m)
	}
	if store.HasSkillMD("alpha-skill") {
		t.Fatal("skill dir should be gone")
	}

	// delete missing
	_, err = del.Call(ctx, args(map[string]any{"name": "alpha-skill"}))
	if err == nil {
		t.Fatal("expected error deleting missing skill")
	}
	if !errors.Is(err, skills.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// delete empty name
	_, err = del.Call(ctx, args(map[string]any{"name": "  "}))
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestCreateSkillTool_allowedToolsAlias(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	store, err := skills.NewStore(root)
	if err != nil {
		t.Fatal(err)
	}
	create := &createSkillTool{store: store}
	raw := json.RawMessage(`{"name":"beta-skill","description":"B","instructions":"I","allowed-tools":"read_web"}`)
	out, err := create.Call(context.Background(), raw)
	if err != nil {
		t.Fatal(err)
	}
	skill := out.(skills.Skill)
	if skill.AllowedTools != "read_web" {
		t.Fatalf("allowed-tools alias: got %q", skill.AllowedTools)
	}
	if !store.HasSkillMD("beta-skill") {
		t.Fatal("expected skill file")
	}
}
