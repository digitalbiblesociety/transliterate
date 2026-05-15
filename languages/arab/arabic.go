// Package arabic transliterates Arabic strings to English.
//
// The default Transliterator is a hybrid: it first looks up the input in a
// dictionary built from the ANETAC corpus (https://github.com/HadjAmeur/ANETAC-Dataset)
// and falls back to a character-level rule-based transliterator for inputs
// not present in the corpus.
package arab

import "sync"

// Source indicates how a Transliterator produced a given result.
type Source int

const (
	// SourceDict means the result came from a dictionary lookup.
	SourceDict Source = iota
	// SourceRules means the result came from the rule-based fallback.
	SourceRules
	// SourcePassthrough means the input contained no Arabic and was returned as-is.
	SourcePassthrough
)

func (s Source) String() string {
	switch s {
	case SourceDict:
		return "dict"
	case SourceRules:
		return "rules"
	case SourcePassthrough:
		return "passthrough"
	}
	return "unknown"
}

// Transliterator converts Arabic text to English.
type Transliterator struct {
	dict *Dictionary
}

// New returns a Transliterator backed by the embedded ANETAC corpus.
func New() (*Transliterator, error) {
	d, err := NewDictionary()
	if err != nil {
		return nil, err
	}
	return &Transliterator{dict: d}, nil
}

// WithDictionary returns a Transliterator backed by the given dictionary.
func WithDictionary(d *Dictionary) *Transliterator {
	return &Transliterator{dict: d}
}

// Transliterate converts the given Arabic string to English and reports the
// source of the result.
func (t *Transliterator) Transliterate(s string) (string, Source) {
	if !containsArabic(s) {
		return s, SourcePassthrough
	}
	if v, ok := t.dict.Lookup(s); ok {
		return v, SourceDict
	}
	return ApplyRules(s), SourceRules
}

// Dictionary exposes the underlying dictionary for direct inspection.
func (t *Transliterator) Dictionary() *Dictionary { return t.dict }

func containsArabic(s string) bool {
	for _, r := range s {
		if (r >= 0x0600 && r <= 0x06FF) || (r >= 0x0750 && r <= 0x077F) {
			return true
		}
	}
	return false
}

var (
	defaultOnce sync.Once
	defaultT    *Transliterator
	defaultErr  error
)

// Default returns a process-wide Transliterator, loading the corpus once.
// Use this from short-lived CLI invocations; library consumers should
// construct their own via New.
func Default() (*Transliterator, error) {
	defaultOnce.Do(func() {
		defaultT, defaultErr = New()
	})
	return defaultT, defaultErr
}

// Transliterate is a convenience that uses the default Transliterator.
func Transliterate(s string) (string, Source, error) {
	t, err := Default()
	if err != nil {
		return "", SourcePassthrough, err
	}
	out, src := t.Transliterate(s)
	return out, src, nil
}
