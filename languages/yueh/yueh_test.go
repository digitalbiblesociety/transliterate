package yueh

import "testing"

func TestJyutpingSamples(t *testing.T) {
	cases := map[string]string{
		"你好": "nei5 hou2",
		"廣州": "gwong2 zau1",
		"一二三": "jat1 ji6 saam1",
	}
	for in, want := range cases {
		if got := Transliterate(in); got != want {
			t.Errorf("Transliterate(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestAtonalStrip(t *testing.T) {
	cases := map[string]string{
		"你好": "nei hou",
		"廣州": "gwong zau",
	}
	for in, want := range cases {
		if got := TransliterateAtonal(in); got != want {
			t.Errorf("TransliterateAtonal(%q) = %q, want %q", in, got, want)
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
