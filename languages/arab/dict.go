package arab

import (
	"bufio"
	"bytes"
	"fmt"
	"io/fs"
	"strings"

	"github.com/digitalbiblesociety/transliterate/internal/data"
)

// Dictionary holds Arabic→English mappings sourced from a parallel corpus.
// Keys are normalized Arabic strings (see normalizeArabic).
type Dictionary struct {
	entries map[string]string
}

// NewDictionary loads every split of the embedded ANETAC corpus.
func NewDictionary() (*Dictionary, error) {
	return LoadDictionary(data.Files, "anetac", data.Splits...)
}

// LoadDictionary loads parallel ar/en files named "<dir>/<split>.ar" and
// "<dir>/<split>.en" from fsys for each requested split.
func LoadDictionary(fsys fs.FS, dir string, splits ...string) (*Dictionary, error) {
	d := &Dictionary{entries: make(map[string]string, 80_000)}
	for _, split := range splits {
		arPath := fmt.Sprintf("%s/%s.ar", dir, split)
		enPath := fmt.Sprintf("%s/%s.en", dir, split)
		if err := d.loadPair(fsys, arPath, enPath); err != nil {
			return nil, fmt.Errorf("load split %q: %w", split, err)
		}
	}
	return d, nil
}

func (d *Dictionary) loadPair(fsys fs.FS, arPath, enPath string) error {
	ar, err := fs.ReadFile(fsys, arPath)
	if err != nil {
		return err
	}
	en, err := fs.ReadFile(fsys, enPath)
	if err != nil {
		return err
	}
	arLines := scanLines(ar)
	enLines := scanLines(en)
	if len(arLines) != len(enLines) {
		return fmt.Errorf("line count mismatch: %s=%d %s=%d", arPath, len(arLines), enPath, len(enLines))
	}
	for i, a := range arLines {
		key := normalizeArabic(a)
		if key == "" {
			continue
		}
		// First occurrence wins; the corpus has a handful of duplicate keys
		// with differing English forms, and we prefer training-set conventions.
		if _, ok := d.entries[key]; !ok {
			d.entries[key] = enLines[i]
		}
	}
	return nil
}

// Lookup returns the English transliteration for an Arabic input and a bool
// indicating whether the entry was found.
func (d *Dictionary) Lookup(arabic string) (string, bool) {
	v, ok := d.entries[normalizeArabic(arabic)]
	return v, ok
}

// Len reports how many unique Arabic keys are loaded.
func (d *Dictionary) Len() int { return len(d.entries) }

func scanLines(b []byte) []string {
	var out []string
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Buffer(make([]byte, 64*1024), 1<<20)
	for s.Scan() {
		out = append(out, strings.TrimSpace(s.Text()))
	}
	return out
}
