package tools

import (
	"context"
	"encoding/json"
	"fmt"
)

// Tool is a callable MCP tool (executable capability). Agent Skills on disk are separate.
type Tool interface {
	Name() string
	// Definition returns one MCP tools/list entry: name, description, inputSchema.
	Definition() map[string]any
	Call(ctx context.Context, arguments json.RawMessage) (structuredContent any, err error)
}

// Registry maps tool names to implementations for tools/list and tools/call.
type Registry struct {
	byName map[string]Tool
	order  []string
}

func NewRegistry() *Registry {
	return &Registry{
		byName: make(map[string]Tool),
	}
}

// Register adds a tool. Duplicate names return an error.
func (r *Registry) Register(t Tool) error {
	name := t.Name()
	if name == "" {
		return fmt.Errorf("tool name is empty")
	}
	if _, exists := r.byName[name]; exists {
		return fmt.Errorf("duplicate tool %q", name)
	}
	r.byName[name] = t
	r.order = append(r.order, name)
	return nil
}

// MustRegister registers or panics (startup only).
func (r *Registry) MustRegister(t Tool) {
	if err := r.Register(t); err != nil {
		panic(err)
	}
}

// Definitions returns MCP tools/list payloads in stable registration order.
func (r *Registry) Definitions() []map[string]any {
	out := make([]map[string]any, 0, len(r.order))
	for _, name := range r.order {
		out = append(out, r.byName[name].Definition())
	}
	return out
}

// Call runs a tool by name and returns structured payload (before MCP envelope).
func (r *Registry) Call(ctx context.Context, name string, arguments json.RawMessage) (any, error) {
	t, ok := r.byName[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool %q", name)
	}
	return t.Call(ctx, arguments)
}
