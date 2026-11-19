package sorensendicecoefficient

import (
	"math"
	"testing"
)

const eps = 1e-9

func almostEqual(t *testing.T, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > eps {
		t.Fatalf("got %v, want %v", got, want)
	}
}

func TestIdenticalStrings(t *testing.T) {
	almostEqual(t, Similarity("night", "night"), 1)
}

func TestCompletelyDifferentStrings(t *testing.T) {
	// "ab" and "cd" share no bigrams, so the coefficient is 0.
	almostEqual(t, Similarity("ab", "cd"), 0)
}

func TestKnownCoefficient(t *testing.T) {
	// "night" bigrams: ni, ig, gh, ht
	// "nacht" bigrams: na, ac, ch, ht
	// Intersection: ht (1). 2*1 / (4+4) = 0.25
	almostEqual(t, Similarity("night", "nacht"), 0.25)
}

func TestRepeatedBigramsCountAsMultiset(t *testing.T) {
	// "aaaa" bigrams: aa, aa, aa
	// "aa"   bigrams: aa
	// Intersection: one aa. 2*1 / (3+1) = 0.5
	// A set-based implementation would wrongly return 2*1/(1+1) = 1.
	almostEqual(t, Similarity("aaaa", "aa"), 0.5)
}

func TestUnicodeRunes(t *testing.T) {
	// "héllo" bigrams: hé, él, ll, lo
	// "hello" bigrams: he, el, ll, lo
	// Intersection: ll, lo (2). 2*2 / (4+4) = 0.5
	almostEqual(t, Similarity("héllo", "hello"), 0.5)
}

func TestEmptyStrings(t *testing.T) {
	almostEqual(t, Similarity("", ""), 1)
	almostEqual(t, Similarity("", "a"), 0)
	almostEqual(t, Similarity("a", ""), 0)
}

func TestSingleRuneStrings(t *testing.T) {
	almostEqual(t, Similarity("a", "a"), 1)
	almostEqual(t, Similarity("a", "b"), 0)
}

func TestOneShortOneLong(t *testing.T) {
	// "a" has no bigrams, "ab" has one. The short string contributes
	// nothing, and the strings differ, so the result is 0.
	almostEqual(t, Similarity("a", "ab"), 0)
}

func TestSymmetric(t *testing.T) {
	got := Similarity("kitten", "sitting")
	want := Similarity("sitting", "kitten")
	if math.Abs(got-want) > eps {
		t.Fatalf("Similarity is not symmetric: %v vs %v", got, want)
	}
}

func TestResultInUnitInterval(t *testing.T) {
	pairs := [][2]string{
		{"", ""},
		{"a", "b"},
		{"ab", "ba"},
		{"aaaa", "aa"},
		{"héllo", "hello"},
		{"kitten", "sitting"},
	}
	for _, p := range pairs {
		got := Similarity(p[0], p[1])
		if got < -eps || got > 1+eps {
			t.Fatalf("Similarity(%q, %q) = %v, outside [0, 1]", p[0], p[1], got)
		}
	}
}
