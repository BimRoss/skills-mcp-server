package skills

import (
	"errors"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func ParseSkillMarkdown(content string, dirName string) (Skill, error) {
	if !strings.HasPrefix(content, "---\n") {
		return Skill{}, errors.New("SKILL.md must start with YAML frontmatter")
	}
	rest := content[len("---\n"):]
	idx := strings.Index(rest, "\n---\n")
	if idx < 0 {
		return Skill{}, errors.New("SKILL.md missing closing frontmatter delimiter")
	}
	frontmatterText := rest[:idx]
	body := strings.TrimSpace(rest[idx+len("\n---\n"):])

	var frontmatter SkillFrontmatter
	if err := yaml.Unmarshal([]byte(frontmatterText), &frontmatter); err != nil {
		return Skill{}, fmt.Errorf("parse frontmatter: %w", err)
	}
	if err := ValidateFrontmatter(frontmatter, dirName); err != nil {
		return Skill{}, err
	}

	return Skill{
		Name:          frontmatter.Name,
		Description:   frontmatter.Description,
		License:       frontmatter.License,
		Compatibility: frontmatter.Compatibility,
		Metadata:      frontmatter.Metadata,
		AllowedTools:  frontmatter.AllowedTools,
		Instructions:  body,
		Raw:           content,
	}, nil
}

func ComposeSkillMarkdown(input CreateOrUpdateSkillInput) (string, error) {
	frontmatter := SkillFrontmatter{
		Name:          input.Name,
		Description:   input.Description,
		License:       input.License,
		Compatibility: input.Compatibility,
		Metadata:      input.Metadata,
		AllowedTools:  input.AllowedTools,
	}
	b, err := yaml.Marshal(frontmatter)
	if err != nil {
		return "", fmt.Errorf("marshal frontmatter: %w", err)
	}
	instructions := strings.TrimSpace(input.Instructions)
	if instructions == "" {
		instructions = "# " + input.Name
	}
	return "---\n" + string(b) + "---\n\n" + instructions + "\n", nil
}
