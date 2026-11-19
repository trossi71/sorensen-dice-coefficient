// Package sorensendicecoefficient computes the Sørensen–Dice coefficient
// between two strings using character bigrams.
//
// The package is deliberately small: it exposes one function, Similarity,
// which returns a value in [0, 1] where 1 means the two strings share every
// bigram and 0 means they share none. The implementation treats a string as a
// sequence of runes, so multi-byte characters are handled correctly without
// any external dependency.
package sorensendicecoefficient

import "strings"

// Similarity returns the Sørensen–Dice coefficient of a and b.
//
// The coefficient is 2*|A∩B| / (|A| + |B|), where A and B are the multisets
// of character bigrams in a and b. We use multisets rather than sets because
// repeated bigrams carry real signal: "aaaa" and "aa" are more similar than
// "abab" and "ab" even though their bigram sets overlap equally.
//
// Edge cases:
//   - If either string has fewer than two runes it contributes no bigrams.
//     When both strings are that short the result is 1 if they are equal and
//     0 otherwise. This keeps the function reflexive without inventing
//     bigrams that do not exist.
//   - The empty string and the empty string are therefore similar (1).
//   - The empty string and any non-empty string are dissimilar (0).
func Similarity(a, b string) float64 {
	bigramsA := bigrams(a)
	bigramsB := bigrams(b)

	if len(bigramsA) == 0 || len(bigramsB) == 0 {
		if a == b {
			return 1
		}
		return 0
	}

	counts := make(map[string]int, len(bigramsA))
	for _, bg := range bigramsA {
		counts[bg]++
	}

	intersection := 0
	for _, bg := range bigramsB {
		if counts[bg] > 0 {
			counts[bg]--
			intersection++
		}
	}

	return 2 * float64(intersection) / float64(len(bigramsA)+len(bigramsB))
}

// bigrams returns the sequence of overlapping two-rune substrings of s.
//
// We build the result as a slice rather than a map so that repeated bigrams
// are preserved; callers that need counts can build them in one pass. The
// function uses strings.Builder to avoid the quadratic allocation pattern of
// repeated string concatenation.
func bigrams(s string) []string {
	runes := []rune(s)
	if len(runes) < 2 {
		return nil
	}

	result := make([]string, 0, len(runes)-1)
	var b strings.Builder
	for i := 0; i < len(runes)-1; i++ {
		b.Reset()
		b.WriteRune(runes[i])
		b.WriteRune(runes[i+1])
		result = append(result, b.String())
	}
	return result
}
