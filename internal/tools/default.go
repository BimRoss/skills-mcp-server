package tools

import (
	"github.com/bimross/skills-mcp-server/internal/readweb"
	"github.com/bimross/skills-mcp-server/internal/skills"
)

// NewDefaultRegistry registers all built-in MCP tools (skill filesystem + read_web).
// Keep tool names and schemas stable; agents-mcp-server and other callers depend on them.
func NewDefaultRegistry(store *skills.Store, readWeb *readweb.Client) *Registry {
	r := NewRegistry()
	r.MustRegister(&listSkillsTool{store: store})
	r.MustRegister(&searchSkillsTool{store: store})
	r.MustRegister(&readSkillTool{store: store})
	r.MustRegister(&listSkillResourcesTool{store: store})
	r.MustRegister(&readSkillResourceTool{store: store})
	r.MustRegister(&getSkillResourceInfoTool{store: store})
	r.MustRegister(&readWebTool{client: readWeb})
	return r
}
