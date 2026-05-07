package googledocs

import "strings"

// DedupeEmails lowercases, trims, and deduplicates addresses.
func DedupeEmails(in []string) []string {
	if len(in) == 0 {
		return nil
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, raw := range in {
		email := strings.ToLower(strings.TrimSpace(raw))
		if email == "" {
			continue
		}
		if _, ok := seen[email]; ok {
			continue
		}
		seen[email] = struct{}{}
		out = append(out, email)
	}
	return out
}

// SubtractEmails returns source emails not present in remove (case-insensitive).
func SubtractEmails(source, remove []string) []string {
	if len(source) == 0 {
		return nil
	}
	rm := map[string]struct{}{}
	for _, raw := range remove {
		key := strings.ToLower(strings.TrimSpace(raw))
		if key != "" {
			rm[key] = struct{}{}
		}
	}
	out := make([]string, 0, len(source))
	for _, raw := range source {
		key := strings.ToLower(strings.TrimSpace(raw))
		if key == "" {
			continue
		}
		if _, exists := rm[key]; exists {
			continue
		}
		out = append(out, key)
	}
	return out
}
