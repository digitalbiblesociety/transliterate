package jpan

import "testing"

// TestKana exercises pure-kana inputs that hit the fast path (no
// tokenizer).
func TestKana(t *testing.T) {
	cases := map[string]string{
		"あいうえお":    "aiueo",
		"きょう":      "kyou", // hiragana digraph; traditional Hepburn keeps the trailing う
		"しょうがっこう": "shougakkou",
		"コーヒー":     "kōhī",     // chōonpu lengthens
		"がっこう":     "gakkou",   // small tsu doubles
		"まっちゃ":     "matcha",   // small tsu before ch → tch
		"しんぶん":     "shimbun",  // n before b → m
		"はんぱ":      "hampa",    // n before p → m
		"たんい":      "tan'i",    // n + vowel → apostrophe
		"ほんや":      "hon'ya",   // n + y → apostrophe
		"":         "",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestKanji exercises kanji and mixed-script inputs that hit kagome.
func TestKanji(t *testing.T) {
	cases := map[string]string{
		"日本":   "nippon",   // common reading
		"漢字":   "kanji",
		"東京":   "toukyou",
		"聖書":   "seisho",   // Bible (聖 + 書)
		"イエス":  "iesu",
		"キリスト": "kirisuto",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestMixedAndPassthrough confirms that USFM markup, ASCII, and digits
// survive untouched alongside Japanese tokens.
func TestMixedAndPassthrough(t *testing.T) {
	cases := map[string]string{
		"Hello":            "Hello",
		"\\v 1 In the":     "\\v 1 In the",
		"Hello 日本":         "Hello nippon",
		"\\v 1 イエスは神の子です": "\\v 1 iesuhakaminokodesu",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestContains(t *testing.T) {
	cases := map[string]bool{
		"":      false,
		"Hello": false,
		"あ":     true,
		"カ":     true,
		"漢":     true,
		"abc 日本 def": true,
	}
	for in, want := range cases {
		if got := Contains(in); got != want {
			t.Errorf("Contains(%q) = %v, want %v", in, got, want)
		}
	}
}
