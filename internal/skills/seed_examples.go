package skills

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// SeedExamplesFromDir copies each direct subdirectory of examplesRoot that looks like an Agent Skill
// (contains SKILL.md) into the store. Skills that already have SKILL.md in the store are skipped
// entirely (idempotent). Missing or non-directory examplesRoot is a no-op.
func SeedExamplesFromDir(s *Store, examplesRoot string) error {
	if examplesRoot == "" {
		return nil
	}
	examplesRoot = filepath.Clean(examplesRoot)
	fi, err := os.Stat(examplesRoot)
	if err != nil || !fi.IsDir() {
		return nil
	}
	entries, err := os.ReadDir(examplesRoot)
	if err != nil {
		return fmt.Errorf("read examples dir: %w", err)
	}
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		name := ent.Name()
		skillSrc := filepath.Join(examplesRoot, name)
		if _, err := os.Stat(filepath.Join(skillSrc, "SKILL.md")); err != nil {
			continue
		}
		if s.HasSkillMD(name) {
			continue
		}
		dest := filepath.Join(s.RootPath(), name)
		if err := copyDirTreeWriteNewOnly(skillSrc, dest); err != nil {
			return fmt.Errorf("seed skill %q: %w", name, err)
		}
	}
	return nil
}

func copyDirTreeWriteNewOnly(srcRoot, destRoot string) error {
	return filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcRoot, path)
		if err != nil {
			return err
		}
		target := filepath.Join(destRoot, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return writeFileIfNew(target, data)
	})
}

func writeFileIfNew(destPath string, data []byte) error {
	if _, err := os.Stat(destPath); err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	tmp := destPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", tmp, err)
	}
	if err := os.Rename(tmp, destPath); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("rename to %s: %w", destPath, err)
	}
	return nil
}
