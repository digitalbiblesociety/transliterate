package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/digitalbiblesociety/transliterate/script"
)

// modeFlags is the single -mode flag plus the three legacy boolean
// aliases. Mode dispatch lives in script.ResolveMode.
type modeFlags struct {
	mode     *string
	tashkeel *bool
	notones  *bool
	phonetic *bool
}

func registerModeFlags(fs *flag.FlagSet) *modeFlags {
	return &modeFlags{
		mode:     fs.String("mode", "", "alternate transliteration mode (run 'translit help' for the list)"),
		tashkeel: fs.Bool("tashkeel", false, "alias for -mode tashkeel (deprecated)"),
		notones:  fs.Bool("notones", false, "alias for -mode atonal (deprecated)"),
		phonetic: fs.Bool("phonetic", false, "alias for -mode phonetic (deprecated)"),
	}
}

func (m *modeFlags) effective() string {
	if *m.mode != "" {
		return *m.mode
	}
	switch {
	case *m.tashkeel:
		return "tashkeel"
	case *m.notones:
		return "atonal"
	case *m.phonetic:
		return "phonetic"
	}
	return ""
}

// validate exits if -mode names a value no script has registered.
// Unknown but irrelevant modes (e.g. -mode tashkeel on Greek input)
// stay silent, matching the old per-flag behaviour.
func (m *modeFlags) validate() {
	name := m.effective()
	if name == "" {
		return
	}
	if slices.Contains(script.ModeNames(), name) {
		return
	}
	fmt.Fprintf(os.Stderr, "unknown mode %q\nvalid: %s\n", name, strings.Join(script.ModeNames(), ", "))
	os.Exit(2)
}

func modesHelp() string {
	var b strings.Builder
	b.WriteString("Modes (-mode <name>):\n")
	for _, m := range script.RegisteredModes() {
		fmt.Fprintf(&b, "  %-9s %s\n", m.Name, m.Description)
	}
	return b.String()
}
