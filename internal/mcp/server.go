package mcp

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bimross/skills-mcp-server/internal/googledocs"
	"github.com/bimross/skills-mcp-server/internal/readweb"
	"github.com/bimross/skills-mcp-server/internal/skills"
	"github.com/bimross/skills-mcp-server/internal/tools"
)

type Server struct {
	registry *tools.Registry
}

func New(store *skills.Store, readWeb *readweb.Client, googleDocs googledocs.EnvConfig) *Server {
	return &Server{registry: tools.NewDefaultRegistry(store, readWeb, googleDocs)}
}

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/mcp", s.handle)
}

type request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type response struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      any         `json:"id,omitempty"`
	Result  any         `json:"result,omitempty"`
	Error   *respErrObj `json:"error,omitempty"`
}

type respErrObj struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeMCP(w, response{JSONRPC: "2.0", Error: &respErrObj{Code: -32700, Message: err.Error()}})
		return
	}

	res := response{JSONRPC: "2.0", ID: req.ID}
	switch req.Method {
	case "initialize":
		res.Result = map[string]any{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]any{
				"name":    "skills-mcp-server",
				"version": "0.1.0",
			},
			"capabilities": map[string]any{
				"tools": map[string]any{},
			},
		}
	case "tools/list":
		res.Result = map[string]any{"tools": s.registry.Definitions()}
	case "tools/call":
		result, err := s.callTool(r.Context(), req.Params)
		if err != nil {
			res.Error = &respErrObj{Code: -32000, Message: err.Error()}
		} else {
			res.Result = result
		}
	default:
		res.Error = &respErrObj{Code: -32601, Message: "method not found"}
	}
	writeMCP(w, res)
}

func (s *Server) callTool(ctx context.Context, raw json.RawMessage) (map[string]any, error) {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(raw, &params); err != nil {
		return nil, err
	}
	structured, err := s.registry.Call(ctx, params.Name, params.Arguments)
	if err != nil {
		return nil, err
	}
	return tools.MCPResult(structured), nil
}

func writeMCP(w http.ResponseWriter, res response) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
