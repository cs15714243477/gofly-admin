package uniapp

import "testing"

func TestNormalizePropertyIDList(t *testing.T) {
	t.Run("comma string with invalid values", func(t *testing.T) {
		got := normalizePropertyIDList("1,2,2,0,-3,abc,5", 20)
		want := []int64{1, 2, 5}
		if len(got) != len(want) {
			t.Fatalf("len mismatch, got=%v want=%v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("idx=%d got=%d want=%d", i, got[i], want[i])
			}
		}
	})

	t.Run("json-like string and max count", func(t *testing.T) {
		got := normalizePropertyIDList("[10,11,12,13]", 2)
		want := []int64{10, 11}
		if len(got) != len(want) {
			t.Fatalf("len mismatch, got=%v want=%v", got, want)
		}
		for i := range want {
			if got[i] != want[i] {
				t.Fatalf("idx=%d got=%d want=%d", i, got[i], want[i])
			}
		}
	})
}

func TestIsPublicDetailView(t *testing.T) {
	trueCases := []string{"1", "true", "TRUE", " yes "}
	for _, c := range trueCases {
		if !isPublicDetailView(c) {
			t.Fatalf("expected true for %q", c)
		}
	}

	falseCases := []string{"", "0", "false", "no", "random"}
	for _, c := range falseCases {
		if isPublicDetailView(c) {
			t.Fatalf("expected false for %q", c)
		}
	}
}
