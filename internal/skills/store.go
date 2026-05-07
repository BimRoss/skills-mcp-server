package skills

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	ErrNotFound         = errors.New("not found")
	ErrInvalidResource  = errors.New("invalid resource path")
	allowedResourceRoot = map[string]struct{}{
		"scripts":    {},
		"references": {},
		"assets":     {},
	}
)

type Store struct {
	root string
}

func NewStore(root string) (*Store, error) {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create skills root: %w", err)
	}
	return &Store{root: root}, nil
}

// RootPath is the filesystem directory backing ListSkills / CRUD.
func (s *Store) RootPath() string {
	return s.root
}

// HasSkillMD reports whether a skill directory with a readable SKILL.md exists.
func (s *Store) HasSkillMD(name string) bool {
	if err := validateName(name); err != nil {
		return false
	}
	p := filepath.Join(s.root, name, "SKILL.md")
	st, err := os.Stat(p)
	return err == nil && !st.IsDir()
}

func (s *Store) ListSkills(query string) ([]Skill, error) {
	entries, err := os.ReadDir(s.root)
	if err != nil {
		return nil, fmt.Errorf("read skills dir: %w", err)
	}
	items := make([]Skill, 0, len(entries))
	q := strings.ToLower(strings.TrimSpace(query))
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		skill, err := s.ReadSkill(e.Name())
		if err != nil {
			continue
		}
		if q != "" {
			text := strings.ToLower(skill.Name + " " + skill.Description + " " + skill.Instructions)
			if !strings.Contains(text, q) {
				continue
			}
		}
		items = append(items, skill)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	return items, nil
}

func (s *Store) ReadSkill(name string) (Skill, error) {
	if err := validateName(name); err != nil {
		return Skill{}, err
	}
	skillPath := filepath.Join(s.root, name, "SKILL.md")
	b, err := os.ReadFile(skillPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Skill{}, ErrNotFound
		}
		return Skill{}, fmt.Errorf("read skill file: %w", err)
	}
	skill, err := ParseSkillMarkdown(string(b), name)
	if err != nil {
		return Skill{}, err
	}
	info, err := os.Stat(skillPath)
	if err == nil {
		skill.UpdatedAt = info.ModTime().UTC()
	}
	return skill, nil
}

func (s *Store) CreateSkill(name string, input CreateOrUpdateSkillInput) (Skill, error) {
	if err := ValidateCreateOrUpdate(input, name); err != nil {
		return Skill{}, err
	}
	dir := filepath.Join(s.root, name)
	if _, err := os.Stat(dir); err == nil {
		return Skill{}, fmt.Errorf("skill %q already exists", name)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return Skill{}, fmt.Errorf("create skill dir: %w", err)
	}
	raw, err := ComposeSkillMarkdown(input)
	if err != nil {
		return Skill{}, err
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), []byte(raw), 0o644); err != nil {
		return Skill{}, fmt.Errorf("write SKILL.md: %w", err)
	}
	return s.ReadSkill(name)
}

func (s *Store) UpdateSkill(name string, input CreateOrUpdateSkillInput) (Skill, error) {
	if _, err := s.ReadSkill(name); err != nil {
		return Skill{}, err
	}
	if err := ValidateCreateOrUpdate(input, name); err != nil {
		return Skill{}, err
	}
	raw, err := ComposeSkillMarkdown(input)
	if err != nil {
		return Skill{}, err
	}
	if err := os.WriteFile(filepath.Join(s.root, name, "SKILL.md"), []byte(raw), 0o644); err != nil {
		return Skill{}, fmt.Errorf("write SKILL.md: %w", err)
	}
	return s.ReadSkill(name)
}

func (s *Store) DeleteSkill(name string) error {
	if err := validateName(name); err != nil {
		return err
	}
	dir := filepath.Join(s.root, name)
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove skill: %w", err)
	}
	return nil
}

func (s *Store) ListResources(name string) ([]ResourceInfo, error) {
	if _, err := s.ReadSkill(name); err != nil {
		return nil, err
	}
	root := filepath.Join(s.root, name)
	items := make([]ResourceInfo, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "SKILL.md" {
			return nil
		}
		rel = filepath.ToSlash(rel)
		first := strings.Split(rel, "/")[0]
		if _, ok := allowedResourceRoot[first]; !ok {
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return nil
		}
		items = append(items, ResourceInfo{
			Path:      rel,
			SizeBytes: info.Size(),
			UpdatedAt: info.ModTime().UTC(),
		})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk resources: %w", err)
	}
	sort.Slice(items, func(i, j int) bool {
		return items[i].Path < items[j].Path
	})
	return items, nil
}

func (s *Store) ReadResource(name, resourcePath string) ([]byte, ResourceInfo, error) {
	fullPath, rel, err := s.resolveResourcePath(name, resourcePath)
	if err != nil {
		return nil, ResourceInfo{}, err
	}
	b, err := os.ReadFile(fullPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, ResourceInfo{}, ErrNotFound
		}
		return nil, ResourceInfo{}, fmt.Errorf("read resource: %w", err)
	}
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, ResourceInfo{}, fmt.Errorf("stat resource: %w", err)
	}
	return b, ResourceInfo{
		Path:      rel,
		SizeBytes: info.Size(),
		UpdatedAt: info.ModTime().UTC(),
	}, nil
}

func (s *Store) WriteResource(name, resourcePath string, content []byte) (ResourceInfo, error) {
	fullPath, rel, err := s.resolveResourcePath(name, resourcePath)
	if err != nil {
		return ResourceInfo{}, err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return ResourceInfo{}, fmt.Errorf("mkdir resource parent: %w", err)
	}
	if err := os.WriteFile(fullPath, content, 0o644); err != nil {
		return ResourceInfo{}, fmt.Errorf("write resource: %w", err)
	}
	_, info, err := s.ReadResource(name, rel)
	return info, err
}

func (s *Store) DeleteResource(name, resourcePath string) error {
	fullPath, _, err := s.resolveResourcePath(name, resourcePath)
	if err != nil {
		return err
	}
	if err := os.Remove(fullPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ErrNotFound
		}
		return fmt.Errorf("delete resource: %w", err)
	}
	return nil
}

func (s *Store) resolveResourcePath(name, resourcePath string) (string, string, error) {
	if _, err := s.ReadSkill(name); err != nil {
		return "", "", err
	}
	clean := filepath.ToSlash(filepath.Clean(resourcePath))
	if clean == "." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return "", "", ErrInvalidResource
	}
	root := strings.Split(clean, "/")[0]
	if _, ok := allowedResourceRoot[root]; !ok {
		return "", "", ErrInvalidResource
	}
	full := filepath.Join(s.root, name, filepath.FromSlash(clean))
	return full, clean, nil
}

func EncodeResourceContent(content []byte, isText bool) string {
	if isText {
		return string(content)
	}
	return base64.StdEncoding.EncodeToString(content)
}

func IsLikelyText(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".txt", ".json", ".yaml", ".yml", ".go", ".sh", ".py", ".js", ".ts":
		return true
	default:
		return false
	}
}

func TimestampNow() time.Time {
	return time.Now().UTC()
}
