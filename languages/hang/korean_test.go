package hang

import "testing"

func TestSamples(t *testing.T) {
	cases := map[string]string{
		"안녕":       "annyeong",    // hello
		"한국":       "hanguk",      // Korea
		"예수":       "yesu",        // Jesus
		"하나님":      "hananim",     // God
		"사랑":       "sarang",      // love
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPassthrough(t *testing.T) {
	for _, s := range []string{"Hello", "\\v 1 In the beginning"} {
		if got := Transliterate(s); got != s {
			t.Errorf("passthrough(%q) = %q", s, got)
		}
	}
}
