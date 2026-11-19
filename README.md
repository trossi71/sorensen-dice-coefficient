# sorensen-dice-coefficient

Computes the Sørensen–Dice coefficient between two strings from their character bigrams.

```go
package main

import (
	"fmt"

	sorensendicecoefficient "sorensen-dice-coefficient"
)

func main() {
	fmt.Println(sorensendicecoefficient.Similarity("night", "nacht")) // 0.25
	fmt.Println(sorensendicecoefficient.Similarity("night", "night")) // 1
}
```

The package exports a single function:

- `Similarity(a, b string) float64` — returns a value in `[0, 1]`.

## Why this exists

String similarity is a common building block for fuzzy matching, and the Sørensen–Dice coefficient over bigrams is a good default because it is cheap, symmetric, and easy to reason about. The trade-off I made is to treat bigrams as a multiset rather than a set: repeated bigrams count more than once. That means `"aaaa"` and `"aa"` score `0.5`, not `1`, which matches the intuition that the shorter string is less similar to the longer one. A set-based implementation would collapse the repetition and overstate the similarity.

## The awkward edge

Strings shorter than two runes have no bigrams. When both inputs are that short the function returns `1` if they are equal and `0` otherwise, so `Similarity("", "")` is `1` and `Similarity("a", "b")` is `0`. This keeps the function reflexive without inventing a bigram that does not exist. If you need a different convention for empty inputs, wrap the call.

## Unicode

Inputs are treated as sequences of runes, so multi-byte characters are handled correctly. `Similarity("héllo", "hello")` is `0.5`.
