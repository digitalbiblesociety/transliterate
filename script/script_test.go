package script

import (
	"reflect"
	"strings"
	"testing"
)

func TestByNameCaseInsensitive(t *testing.T) {
	cases := []string{"Thai", "thai", "THAI", "tHaI"}
	for _, name := range cases {
		eng := ByName(name)
		if eng == nil {
			t.Fatalf("ByName(%q): nil", name)
		}
		if eng.Name != "Thai" {
			t.Errorf("ByName(%q): got %s, want Thai", name, eng.Name)
		}
	}
	if ByName("Latn") != nil {
		t.Errorf("ByName(\"Latn\") = non-nil; Latin isn't registered")
	}
}

func TestDetectThai(t *testing.T) {
	eng := Detect("ในปฐมกาล")
	if eng == nil {
		t.Fatal("Detect returned nil for Thai input")
	}
	if eng.Name != "Thai" {
		t.Errorf("Detect Thai: got %s, want Thai", eng.Name)
	}
}

func TestDetectGreek(t *testing.T) {
	eng := Detect("Ἰησοῦς Χριστός")
	if eng == nil {
		t.Fatal("Detect returned nil for Greek input")
	}
	if eng.Name != "Grek" {
		t.Errorf("Detect Greek: got %s, want Grek", eng.Name)
	}
}

func TestDetectPureLatinReturnsNil(t *testing.T) {
	if eng := Detect("Hello, world."); eng != nil {
		t.Errorf("Detect(\"Hello, world.\") = %s; want nil", eng.Name)
	}
}

func TestDetectMixedPicksDominant(t *testing.T) {
	// Greek dominates over a stray Cyrillic letter — Detect should
	// return Grek, not Cyrl.
	eng := Detect("Ἰησοῦς Χριστός Я")
	if eng == nil || eng.Name != "Grek" {
		t.Errorf("Detect(mixed gr+cy): got %v, want Grek", eng)
	}
}

func TestThaiTransliterateRoundtripSmoke(t *testing.T) {
	eng := ByName("Thai")
	if eng == nil {
		t.Fatal("Thai engine missing")
	}
	got := eng.Transliterate("สวัสดี")
	if got == "" || got == "สวัสดี" {
		t.Errorf("Thai Transliterate noop: got %q", got)
	}
	if strings.ContainsAny(got, "฀-๿") {
		t.Errorf("Thai output still contains Thai runes: %q", got)
	}
}

func TestThaiSplitUsesNativeSplitter(t *testing.T) {
	eng := ByName("Thai")
	// Thai splitter splits on spaces (phrase boundaries) and the ฯ
	// pause marker. SplitGeneric would not split on regular spaces.
	got := eng.Split("ก ข ค")
	want := []string{"ก", "ข", "ค"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Thai Split(spaces): got %q, want %q", got, want)
	}
	got = eng.Split("กฯข")
	want = []string{"ก", "ข"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Thai Split(ฯ): got %q, want %q", got, want)
	}
}

func TestNonThaiSplitUsesGeneric(t *testing.T) {
	eng := ByName("Cyrl")
	// SplitGeneric does NOT split on regular spaces.
	got := eng.Split("Москва Київ")
	want := []string{"Москва Київ"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Cyrl Split(spaces): got %q, want %q", got, want)
	}
	got = eng.Split("Москва, Київ")
	want = []string{"Москва", "Київ"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Cyrl Split(comma): got %q, want %q", got, want)
	}
}

func TestNamesIncludesEveryRegistered(t *testing.T) {
	names := Names()
	for _, want := range []string{"Thai", "Grek", "Hebr", "Arab", "Cyrl"} {
		found := false
		for _, n := range names {
			if n == want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Names() missing %q; got %v", want, names)
		}
	}
}

func TestEveryEngineHasAllFunctions(t *testing.T) {
	for _, e := range Engines() {
		if e.Name == "" {
			t.Errorf("engine with empty Name")
		}
		if e.Contains == nil {
			t.Errorf("%s: nil Contains", e.Name)
		}
		if e.Transliterate == nil {
			t.Errorf("%s: nil Transliterate", e.Name)
		}
		if e.Split == nil {
			t.Errorf("%s: nil Split", e.Name)
		}
	}
}
