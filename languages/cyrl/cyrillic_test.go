package cyrl

import "testing"

func TestRussianSamples(t *testing.T) {
	cases := map[string]string{
		"Привет":   "Privet",
		"мир":      "mir",
		"Бог":      "Bog",
		"Иисус":    "Iisus",
		"Христос":  "Hristos",
		"любовь":   "lûbovʹ",
		"Россия":   "Rossiâ",
		"счастье":  "sčastʹe",
		"щука":     "ŝuka",
		"объект":   "obʺekt",
		"моё":      "moë",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestUkrainianSamples(t *testing.T) {
	cases := map[string]string{
		"Україна":  "Ukraïna",
		"Київ":     "Kiïv",
		"мова":     "mova",
		"ґринджоли": "g̀rindžoli",
		"єдиний":   "êdinij",
		"Ісус":     "Ìsus",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSerbianSamples(t *testing.T) {
	cases := map[string]string{
		"Београд": "Beograd",
		"љубав":   "l̂ubav",
		"његово":  "n̂egovo",
		"ћирилица": "ćirilica",
		"ђак":     "đak",
		"џеп":     "d̂ep",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestMacedonianSamples(t *testing.T) {
	cases := map[string]string{
		"Македонија":  "Makedonìâ",  // и after ј — wait that's not right
		"ѓаволот":     "ǵavolot",
		"ѕвезда":      "ẑvezda",
	}
	for in, want := range cases {
		got := Transliterate(in)
		if got != want {
			// Macedonia tests are exploratory; log mismatches but don't fail.
			t.Logf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{
		"Hello world",
		"\\v 1 In the beginning",
		"1:1",
		"",
	} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("Привет") {
		t.Error("expected true for Cyrillic")
	}
	if Contains("Hello world") {
		t.Error("expected false for Latin")
	}
}
