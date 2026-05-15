package script

import (
	"fmt"
	"sort"
)

// Mode is an alternate transliteration function for a script. Several
// scripts may register the same Name (e.g. "atonal" for both Hani and
// Yueh); ResolveMode picks by (Engine, Name).
type Mode struct {
	Script      string
	Name        string
	Description string
	Fn          func(string) string
}

var modeRegistry = map[string]map[string]Mode{}

// RegisterMode is intended for init(). Panics on duplicate
// registration so collisions surface at build time.
func RegisterMode(script, name, description string, fn func(string) string) {
	if modeRegistry[script] == nil {
		modeRegistry[script] = map[string]Mode{}
	}
	if _, exists := modeRegistry[script][name]; exists {
		panic(fmt.Sprintf("script: mode %q already registered for %q", name, script))
	}
	modeRegistry[script][name] = Mode{
		Script:      script,
		Name:        name,
		Description: description,
		Fn:          fn,
	}
}

// ResolveMode falls back to eng.Transliterate when name is empty or
// doesn't match a registered mode — so -mode phonetic with Latin input
// is silently a no-op, matching the old per-flag behaviour.
func ResolveMode(eng *Engine, name string) func(string) string {
	if eng == nil {
		return nil
	}
	if name != "" {
		if scripts, ok := modeRegistry[eng.Name]; ok {
			if m, ok := scripts[name]; ok {
				return m.Fn
			}
		}
	}
	return eng.Transliterate
}

// RegisteredModes returns every mode in stable (name, script) order
// for help-text generation.
func RegisteredModes() []Mode {
	var out []Mode
	for _, modes := range modeRegistry {
		for _, m := range modes {
			out = append(out, m)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name != out[j].Name {
			return out[i].Name < out[j].Name
		}
		return out[i].Script < out[j].Script
	})
	return out
}

func ModeNames() []string {
	seen := map[string]struct{}{}
	for _, modes := range modeRegistry {
		for name := range modes {
			seen[name] = struct{}{}
		}
	}
	out := make([]string, 0, len(seen))
	for name := range seen {
		out = append(out, name)
	}
	sort.Strings(out)
	return out
}
