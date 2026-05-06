package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/bimross/skills-mcp-server/internal/readweb"
	"github.com/bimross/skills-mcp-server/internal/skills"
)

type Handler struct {
	store   *skills.Store
	readWeb *readweb.Client
}

func New(store *skills.Store, readWeb *readweb.Client) *Handler {
	return &Handler{store: store, readWeb: readWeb}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.health)
	mux.HandleFunc("/api/runtime/read-web", h.runtimeReadWeb)
	mux.HandleFunc("/api/skills", h.skillsCollection)
	mux.HandleFunc("/api/skills/", h.skillsResource)
}

func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) runtimeReadWeb(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Query string `json:"query"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	result, err := h.readWeb.Run(r.Context(), input.Query)
	if err != nil {
		writeErr(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"fallbackText": result.Summary,
		"finalSummary": result.Summary,
		"citations":    result.Citations,
	})
}

func (h *Handler) skillsCollection(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := h.store.ListSkills(r.URL.Query().Get("q"))
		if err != nil {
			writeErr(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"skills": items})
	case http.MethodPost:
		var input skills.CreateOrUpdateSkillInput
		if err := decodeJSON(r, &input); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		item, err := h.store.CreateSkill(input.Name, input)
		if err != nil {
			writeStatusFromErr(w, err)
			return
		}
		writeJSON(w, http.StatusCreated, item)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *Handler) skillsResource(w http.ResponseWriter, r *http.Request) {
	trimmed := strings.TrimPrefix(r.URL.Path, "/api/skills/")
	parts := strings.Split(trimmed, "/")
	if len(parts) == 0 || parts[0] == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	name := parts[0]
	if len(parts) == 1 {
		h.singleSkill(w, r, name)
		return
	}
	if parts[1] != "resources" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	h.skillResources(w, r, name, strings.Join(parts[2:], "/"))
}

func (h *Handler) singleSkill(w http.ResponseWriter, r *http.Request, name string) {
	switch r.Method {
	case http.MethodGet:
		item, err := h.store.ReadSkill(name)
		if err != nil {
			writeStatusFromErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodPut:
		var input skills.CreateOrUpdateSkillInput
		if err := decodeJSON(r, &input); err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		item, err := h.store.UpdateSkill(name, input)
		if err != nil {
			writeStatusFromErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, item)
	case http.MethodDelete:
		if err := h.store.DeleteSkill(name); err != nil {
			writeStatusFromErr(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

type writeResourceRequest struct {
	Path            string `json:"path"`
	Content         string `json:"content"`
	ContentBase64   string `json:"contentBase64,omitempty"`
	ContentEncoding string `json:"contentEncoding,omitempty"`
}

func (h *Handler) skillResources(w http.ResponseWriter, r *http.Request, name, resourcePath string) {
	if resourcePath == "" {
		if r.Method == http.MethodGet {
			items, err := h.store.ListResources(name)
			if err != nil {
				writeStatusFromErr(w, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"resources": items})
			return
		}
		if r.Method == http.MethodPost {
			var input writeResourceRequest
			if err := decodeJSON(r, &input); err != nil {
				writeErr(w, http.StatusBadRequest, err)
				return
			}
			content := []byte(input.Content)
			if input.ContentEncoding == "base64" {
				decoded, err := base64.StdEncoding.DecodeString(input.ContentBase64)
				if err != nil {
					writeErr(w, http.StatusBadRequest, err)
					return
				}
				content = decoded
			}
			info, err := h.store.WriteResource(name, input.Path, content)
			if err != nil {
				writeStatusFromErr(w, err)
				return
			}
			writeJSON(w, http.StatusCreated, info)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		content, info, err := h.store.ReadResource(name, resourcePath)
		if err != nil {
			writeStatusFromErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"resource": map[string]any{
				"path":           info.Path,
				"sizeBytes":      info.SizeBytes,
				"updatedAt":      info.UpdatedAt,
				"isText":         skills.IsLikelyText(info.Path),
				"contentEncoded": skills.EncodeResourceContent(content, skills.IsLikelyText(info.Path)),
				"encoding":       textOrBase64(skills.IsLikelyText(info.Path)),
			},
		})
	case http.MethodPut:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			writeErr(w, http.StatusBadRequest, err)
			return
		}
		info, err := h.store.WriteResource(name, resourcePath, body)
		if err != nil {
			writeStatusFromErr(w, err)
			return
		}
		writeJSON(w, http.StatusOK, info)
	case http.MethodDelete:
		if err := h.store.DeleteResource(name, resourcePath); err != nil {
			writeStatusFromErr(w, err)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func textOrBase64(isText bool) string {
	if isText {
		return "text"
	}
	return "base64"
}

func decodeJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dst)
}

func writeStatusFromErr(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, skills.ErrNotFound):
		writeErr(w, http.StatusNotFound, err)
	case errors.Is(err, skills.ErrInvalidResource):
		writeErr(w, http.StatusBadRequest, err)
	default:
		writeErr(w, http.StatusBadRequest, err)
	}
}

func writeErr(w http.ResponseWriter, code int, err error) {
	writeJSON(w, code, map[string]string{"error": err.Error()})
}

func writeJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(payload)
}
