package hani

import "testing"

func TestTonalSamples(t *testing.T) {
	cases := map[string]string{
		"你好":  "nǐ hǎo",
		"中国":  "zhōng guó",
		"小":   "xiǎo",
		"一二三": "yī èr sān",
		"上帝":  "shàng dì",   // God
		"耶稣":  "yē sū",      // Jesus
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestHeteronymDisambiguation exercises characters whose reading
// changes by context. These all rely on the phrase dictionary —
// per-character Unihan lookups would emit the wrong primary reading.
func TestHeteronymDisambiguation(t *testing.T) {
	cases := map[string]string{
		"中国":   "zhōng guó", // 中 = zhōng
		"击中":   "jī zhòng",  // 中 = zhòng
		"行业":   "háng yè",   // 行 = háng
		"行走":   "xíng zǒu",  // 行 = xíng
		"朝阳":   "zhāo yáng", // 朝 = zhāo (the first-listed reading)
		"重要":   "zhòng yào", // 重 = zhòng
		"重新":   "chóng xīn", // 重 = chóng
		"长城":   "cháng chéng", // 长 = cháng
		"长大":   "zhǎng dà",    // 长 = zhǎng
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAtonalStrip(t *testing.T) {
	cases := map[string]string{
		"你好":  "ni hao",
		"中国":  "zhong guo",
		"绿":   "lu", // 绿 → lǜ → lu (ü with tone stripped)
	}
	for in, want := range cases {
		if got := TransliterateAtonal(in); got != want {
			t.Errorf("TransliterateAtonal(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "1:1", "", "你 (parenthesis)"} {
		got := Transliterate(s)
		if !Contains(s) {
			if got != s {
				t.Errorf("non-Han %q: got %q, want passthrough", s, got)
			}
		}
	}
}

func TestContains(t *testing.T) {
	if !Contains("你好") {
		t.Error("expected true for Han string")
	}
	if Contains("Hello") {
		t.Error("expected false for Latin string")
	}
}
