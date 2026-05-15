package arab

import "testing"

// Concrete-case tests for the headline examples we used to design the
// scheme.
func TestTashkeelHeadlineCases(t *testing.T) {
	cases := map[string]string{
		"عَائِلَةُ":   "eayilat", // family — Google's example
		"يَعْقُوبَ":  "yaequb",  // Jacob
		"فِي":          "fi",
		"مِصْر":        "misr",
		"بَرَكَةٍ":     "barakat", // ة+vowel→t; Google itself is inconsistent here
		"شَابٌّ":       "shabb",
		"رَائِحَةُ":   "rayihat",
		"أغْصَانٍ":    "aghsan",
		"ابْنَتُهُ":  "abnatuh",
		"الحَاكِمِ":    "alhakim",
		"حَارَانَ":     "haran",
	}
	for in, want := range cases {
		if got := TransliterateTashkeel(in); got != want {
			t.Errorf("TransliterateTashkeel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestTashkeelPassthrough(t *testing.T) {
	for _, s := range []string{"Hello world", "\\v 1 In the beginning", "1:1"} {
		if got := TransliterateTashkeel(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}
