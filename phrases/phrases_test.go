package phrases

import (
	"reflect"
	"testing"
)

func TestSplitGeneric(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want []string
	}{
		{"empty", "", nil},
		{"single fragment", "hello", []string{"hello"}},
		{"inner spaces are kept", "hello world", []string{"hello world"}},
		{"comma boundary", "alpha,beta", []string{"alpha", "beta"}},
		{"period and space", "first. second.", []string{"first", "second"}},
		{"all boundaries yield nothing", ".\t,\n;", nil},
		{"leading/trailing whitespace trimmed", "  hi  ", []string{"hi"}},
		{
			"newline and tab split, comma splits",
			"one, two\nthree;four",
			[]string{"one", "two", "three", "four"},
		},
		{
			"multi-word phrases preserved",
			"Yes! Or no? Maybe so.",
			[]string{"Yes", "Or no", "Maybe so"},
		},
		{
			"non-Latin passes through verbatim",
			"Москва, Київ; София",
			[]string{"Москва", "Київ", "София"},
		},
		{
			"tabs split, regular spaces don't",
			"col1\tcol2 with words\tcol3",
			[]string{"col1", "col2 with words", "col3"},
		},
		{
			"line separator U+2028 splits",
			"line one line two",
			[]string{"line one", "line two"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SplitGeneric(tc.in)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("SplitGeneric(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestPair(t *testing.T) {
	orig := []string{"alpha", "beta", "gamma"}
	upper := func(s string) string {
		r := []rune(s)
		if len(r) > 0 && r[0] >= 'a' && r[0] <= 'z' {
			r[0] -= 32
		}
		return string(r)
	}
	got := Pair(orig, upper)
	want := []Phrase{
		{ID: "f0001", Original: "alpha", Latin: "Alpha"},
		{ID: "f0002", Original: "beta", Latin: "Beta"},
		{ID: "f0003", Original: "gamma", Latin: "Gamma"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Pair: got %+v, want %+v", got, want)
	}
}

func TestPairEmpty(t *testing.T) {
	got := Pair(nil, func(s string) string { return s })
	if len(got) != 0 {
		t.Errorf("Pair(nil, …): expected empty, got %+v", got)
	}
}

// IDs must be zero-padded to 4 digits so lexicographic and numeric
// orderings match — important for Aeneas sync-map joins.
func TestPairIDFormat(t *testing.T) {
	orig := make([]string, 12)
	for i := range orig {
		orig[i] = "x"
	}
	got := Pair(orig, func(s string) string { return s })
	checks := map[int]string{
		0:  "f0001",
		9:  "f0010",
		11: "f0012",
	}
	for idx, want := range checks {
		if got[idx].ID != want {
			t.Errorf("got[%d].ID = %q, want %q", idx, got[idx].ID, want)
		}
	}
}
