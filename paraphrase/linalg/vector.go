package linalg

import "math"

// a sparse vector of uint64->float64s
type IFVector map[uint64]float64
type IFVectorEach func(key uint64, value float64)
type IFVectorReduce func(key uint64, value, lastValue float64) float64
type IFVectorMap func(key uint64, value float64) float64
type IFVectorFilter func(key uint64, value float64) bool

func NewIfVector() IFVector {
	return make(IFVector)
}

func (vec IFVector) ForEach(apply IFVectorEach) {
	for key, val := range vec {
		apply(key, val)
	}
}

func (vec IFVector) Reduce(init float64, reducer IFVectorReduce) float64 {
	for key, val := range vec {
		init = reducer(key, val, init)
	}

	return init
}

func (vec IFVector) Map(mapper IFVectorMap) IFVector {
	for key, val := range vec {
		vec[key] = mapper(key, val)
	}

	return vec
}

func (vec IFVector) Filter(filter IFVectorFilter) {
	for key, val := range vec {
		keep := filter(key, val)
		if !keep {
			delete(vec, key)
		}
	}
}

func (vec IFVector) FilterMap(filter IFVectorFilter, mapper IFVectorMap) {
	for key, val := range vec {
		keep := filter(key, val)
		if keep {
			vec[key] = mapper(key, val)
		} else {
			delete(vec, key)
		}
	}
}

func (vec IFVector) Sum() float64 {
	return vec.Reduce(0, func(k uint64, v, last float64) float64 { return last + v })
}

func (vec IFVector) Max() float64 {
	max := func(k uint64, v, last float64) float64 {
		return math.Max(v, last)
	}

	return vec.Reduce(0, max)
}

func (vec IFVector) L1Norm() float64 {
	l1norm := func(key uint64, curr, last float64) float64 {
		return math.Abs(curr) + last
	}

	return vec.Reduce(0, l1norm)
}

func (vec IFVector) L2Norm() float64 {
	l2norm := func(key uint64, curr, last float64) float64 {
		return (curr * curr) + last
	}

	return math.Sqrt(vec.Reduce(0, l2norm))
}

func (vec IFVector) Clone() IFVector {
	clone := NewIfVector()

	vec.ForEach(func(key uint64, value float64) { clone[key] = value })

	return clone
}

// Dot computes the dot product of the two vectors
func (vec IFVector) Dot(multiplier IFVector) float64 {
	var sum float64

	small, large := orderIFVectors(vec, multiplier)

	small.ForEach(func(key uint64, value float64) {
		sum += large[key] * value
	})

	return sum
}

func (vec IFVector) MultF(multiplier float64) {
	mult := func(key uint64, val float64) float64 {
		return multiplier * val
	}

	vec.Map(mult)
}

// Prod is an element-wise product between vectors.
// a.Prod(b) = [a[1] * b[1], a[2] * b[2], ...]
func (vec IFVector) Prod(multiplier IFVector) {
	mult := func(key uint64, val float64) float64 {
		return multiplier[key] * val
	}

	vec.Map(mult)
}

// DivF divides each element in the vector by the given value.
func (vec IFVector) DivF(divisor float64) {
	div := func(key uint64, dividend float64) float64 {
		return dividend / divisor
	}

	vec.Map(div)
}

func (vec IFVector) CosineSimilarity(other IFVector) float64 {
	x, y := orderIFVectors(vec, other)

	return x.Dot(y) / (x.L2Norm() * y.L2Norm())
}

func orderIFVectors(a, b IFVector) (smaller, larger IFVector) {
	if len(a) < len(b) {
		return a, b
	} else {
		return b, a
	}
}
