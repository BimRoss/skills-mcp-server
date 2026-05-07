package skills

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	namePattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)
)

func ValidateFrontmatter(frontmatter SkillFrontmatter, dirName string) error {
	if err := validateName(frontmatter.Name); err != nil {
		return err
	}
	if frontmatter.Name != dirName {
		return fmt.Errorf("name %q must match parent directory %q", frontmatter.Name, dirName)
	}
	if len(frontmatter.Description) == 0 || len(frontmatter.Description) > 1024 {
		return errors.New("description must be 1-1024 characters")
	}
	if frontmatter.Compatibility != "" && (len(frontmatter.Compatibility) < 1 || len(frontmatter.Compatibility) > 500) {
		return errors.New("compatibility must be 1-500 characters when provided")
	}
	return nil
}

func ValidateCreateOrUpdate(input CreateOrUpdateSkillInput, pathName string) error {
	if strings.TrimSpace(input.Name) == "" {
		return errors.New("name is required")
	}
	if input.Name != pathName {
		return fmt.Errorf("name %q must match path %q", input.Name, pathName)
	}
	if err := validateName(input.Name); err != nil {
		return err
	}
	if len(strings.TrimSpace(input.Description)) == 0 || len(input.Description) > 1024 {
		return errors.New("description must be 1-1024 characters")
	}
	if input.Compatibility != "" && len(input.Compatibility) > 500 {
		return errors.New("compatibility must be <=500 characters")
	}
	return nil
}

func validateName(name string) error {
	if len(name) < 1 || len(name) > 64 {
		return errors.New("name must be 1-64 characters")
	}
	if !namePattern.MatchString(name) {
		return errors.New("name must use lowercase letters, numbers, and single hyphens only")
	}
	if strings.Contains(name, "--") {
		return errors.New("name must not contain consecutive hyphens")
	}
	return nil
}
