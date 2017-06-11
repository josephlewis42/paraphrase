package paraphrase

import (
	"hash/fnv"
	"math"
	"regexp"
)

var (
	whitespace = regexp.MustCompile(`\s`)
)

type Fingerprint uint64

func normalizeDocument(document []byte) []byte {
	return removeWhitespace(document)
}

func removeWhitespace(document []byte) []byte {
	// https://github.com/golang/go/wiki/SliceTricks
	output := make([]byte, 0, len(document))
	for _, x := range document {

		switch x {
		case '\t', '\n', '\v', '\f', '\r', ' ', 0x85, 0xA0:
			continue
		default:
			output = append(output, x)
		}
	}
	return output
}

func fingerprintDocument(document []byte, size int) []Fingerprint {
	fingerprintCount := len(document) - size

	var fingerprints []Fingerprint

	for i := 0; i <= fingerprintCount; i++ {
		hash := fnv.New64()
		hash.Write(document[i : i+size])
		fingerprints = append(fingerprints, Fingerprint(hash.Sum64()))
	}

	return fingerprints
}

func winnow(fingerprints []Fingerprint, window int, robust bool) []Fingerprint {
	var recorded []Fingerprint

	h := make([]Fingerprint, window)

	for i, _ := range h {
		h[i] = math.MaxUint64
	}

	r := 0   // window right end
	min := 0 // index of min hash

	for _, fingerprint := range fingerprints {
		r = (r + 1) % window // shift window by one
		h[r] = fingerprint

		if min == r {
			// previous minimum is no longer in this window.
			// scan h leftward starting from r for the rightmost minimal hash.
			// Note min starts with the index of the rightmost hash.
			for i := (r - 1 + window) % window; i != r; i = (i - 1 + window) % window {
				if h[i] < h[min] {
					min = i
				}
			}

			recorded = append(recorded, fingerprint)

		} else {
			// Otherwise, the previous minimum is still in this window. Compare
			// against the new value and update min if necessary.
			if h[r] < h[min] || (!robust && h[r] == h[min]) {
				min = r
				recorded = append(recorded, fingerprint)
			}
		}
	}

	return recorded
}

func (p *ParaphraseDb) WinnowData(bytes []byte) (TermCountVector, error) {
	winnowed := make(TermCountVector)

	norm := normalizeDocument(bytes)
	prints := fingerprintDocument(norm, p.settings.FingerprintSize)
	saved := winnow(prints, p.settings.WindowSize, p.settings.RobustHash)

	for _, print := range saved {
		curr := winnowed[uint64(print)]
		winnowed[uint64(print)] = curr + 1
	}

	return winnowed, nil
}
