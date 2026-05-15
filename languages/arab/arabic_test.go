package arab

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/digitalbiblesociety/transliterate/internal/data"
)

func TestDictionaryLoads(t *testing.T) {
	d, err := NewDictionary()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	// Normalization (alef variants, ta marbuta, etc.) collapses some keys,
	// so the unique-key count is lower than the 79,921 raw pairs.
	if d.Len() < 60_000 {
		t.Fatalf("expected >=60k unique keys, got %d", d.Len())
	}
}

func TestLookupKnownEntries(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// Inputs known to be in the corpus; we only assert source=dict and
	// non-empty output, not a specific spelling — multiple English forms
	// can map to the same Arabic name, and train-set conventions win.
	for _, in := range []string{"باراندياران", "كاراسين", "اريكا", "محمد"} {
		got, src := tr.Transliterate(in)
		if src != SourceDict {
			t.Errorf("%q: source = %s, want dict", in, src)
		}
		if got == "" {
			t.Errorf("%q: empty result", in)
		}
	}
}

func TestRulesFallback(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatal(err)
	}
	// Unlikely to be in the corpus; should fall through to rules.
	got, src := tr.Transliterate("xyzابتث")
	if src == SourceDict {
		t.Fatalf("expected rules fallback, got dict result %q", got)
	}
}

func TestPassthroughForNonArabic(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatal(err)
	}
	got, src := tr.Transliterate("Hello World")
	if src != SourcePassthrough {
		t.Errorf("source = %s, want passthrough", src)
	}
	if got != "Hello World" {
		t.Errorf("got %q, want passthrough", got)
	}
}

func TestNormalizeFoldsAlefVariants(t *testing.T) {
	for _, in := range []string{"أحمد", "إحمد", "آحمد", "ٱحمد"} {
		if got := normalizeArabic(in); got != "احمد" {
			t.Errorf("normalize(%q) = %q, want \"احمد\"", in, got)
		}
	}
}

// TestCorpusDictionaryAccuracy measures hit rate on the dev split. The dev
// pairs are part of the embedded corpus, so this confirms the loader sees
// every line and that normalization doesn't drop valid keys.
func TestCorpusDictionaryAccuracy(t *testing.T) {
	tr, err := New()
	if err != nil {
		t.Fatal(err)
	}
	arBytes, err := data.Files.ReadFile("anetac/dev.ar")
	if err != nil {
		t.Fatal(err)
	}
	enBytes, err := data.Files.ReadFile("anetac/dev.en")
	if err != nil {
		t.Fatal(err)
	}
	ar := scan(arBytes)
	en := scan(enBytes)
	if len(ar) != len(en) {
		t.Fatalf("dev mismatch %d vs %d", len(ar), len(en))
	}
	var dictHits, exactMatches int
	for i, a := range ar {
		out, src := tr.Transliterate(a)
		if src == SourceDict {
			dictHits++
		}
		if out == en[i] {
			exactMatches++
		}
	}
	if dictHits != len(ar) {
		t.Errorf("dict hits = %d/%d; expected full coverage of dev set", dictHits, len(ar))
	}
	t.Logf("dev exact-match rate: %d/%d (%.1f%%)", exactMatches, len(ar), 100*float64(exactMatches)/float64(len(ar)))
}

func scan(b []byte) []string {
	var out []string
	s := bufio.NewScanner(bytes.NewReader(b))
	s.Buffer(make([]byte, 64*1024), 1<<20)
	for s.Scan() {
		out = append(out, s.Text())
	}
	return out
}
