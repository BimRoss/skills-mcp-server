package skills

import "testing"

func TestValidateName(t *testing.T) {
	valid := []string{"read-web", "a1", "joanne-tools"}
	for _, name := range valid {
		if err := validateName(name); err != nil {
			t.Fatalf("expected valid name %q: %v", name, err)
		}
	}

	invalid := []string{"Read-Web", "-foo", "foo-", "foo--bar", "foo_bar"}
	for _, name := range invalid {
		if err := validateName(name); err == nil {
			t.Fatalf("expected invalid name %q", name)
		}
	}
}
