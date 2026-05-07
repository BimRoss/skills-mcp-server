package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSeedExamplesFromDir_Idempotent(t *testing.T) {
	t.Parallel()
	exRoot := filepath.Join("..", "..", "examples")
	if _, err := os.Stat(filepath.Join(exRoot, "read-web", "SKILL.md")); err != nil {
		t.Skip("examples/read-web not present at repo root")
	}
	dir := t.TempDir()
	store, err := NewStore(dir)
	if err != nil {
		t.Fatal(err)
	}
	if err := SeedExamplesFromDir(store, exRoot); err != nil {
		t.Fatal(err)
	}
	if !store.HasSkillMD("read-web") {
		t.Fatal("expected read-web seeded")
	}
	if err := SeedExamplesFromDir(store, exRoot); err != nil {
		t.Fatal(err)
	}
}
