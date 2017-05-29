package paraphrase

type BucketSet map[uint64]float64

func NewBucketSet() BucketSet {
	return make(BucketSet)
}

func (hs BucketSet) AddAll(elements []uint64) BucketSet {

	for _, element := range elements {
		val, ok := hs[element]

		if !ok {
			hs[element] = 1
		} else {
			hs[element] = val + 1
		}
	}

	return hs
}

func (hs BucketSet) Intersect(other BucketSet) BucketSet {
	intersection := NewBucketSet()

	for key, val := range hs {
		if _, ok := other[key]; ok {
			intersection[key] = val
		}
	}

	return intersection
}

func (hs BucketSet) Union(other BucketSet) BucketSet {
	union := NewBucketSet()

	for key, val := range other {
		union[key] = val
	}

	for key, val := range hs {
		union[key] = val
	}

	return union
}

func (hs BucketSet) GetOrDefault(key uint64, defaultVal float64) float64 {
	val, ok := hs[key]

	if ok {
		return val
	}

	return defaultVal
}

func (hs BucketSet) Mult(other BucketSet) BucketSet {
	output := hs.Union(other)

	for key, _ := range output {
		output[key] = hs.GetOrDefault(key, 0.0) * other.GetOrDefault(key, 0.0)
	}

	return output
}

func (hs BucketSet) Sum() float64 {
	sum := 0.0

	for _, val := range hs {
		sum += val
	}

	return sum
}

func (hs BucketSet) TfIdf() float64 {
	var tfidf float64
	tfidf = 1.0

	for _, elemFreq := range hs {
		tfidf *= 1.0 / float64(elemFreq)
	}

	return tfidf
}

func (hs BucketSet) OverlapProportion(other BucketSet) float64 {
	return float64(len(hs.Union(other))) / float64(len(hs))
}
