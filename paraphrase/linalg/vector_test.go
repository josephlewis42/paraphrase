package linalg

import (
	"fmt"
	"testing"
)

var (
	result64 float64
)

func ExampleNewIfVector() {
	// create a sparse vector
	vec := NewIfVector()
	fmt.Println(len(vec))

	// add a value
	vec[100] = 1.06
	fmt.Println(len(vec))

	// try getting a non-existant value
	_, ok := vec[0]
	fmt.Println(ok)

	val := vec[100]
	fmt.Println(val)

	// Output:
	// 0
	// 1
	// false
	// 1.06
}

func ExampleIFVector_ForEach() {
	vec := NewIfVector()
	vec[1] = 2
	vec[2] = 4
	vec[3] = 8
	vec[4] = 16

	var keysum uint64
	var valsum float64

	vec.ForEach(func(key uint64, val float64) {
		keysum += key
		valsum += val
	})

	fmt.Println(keysum)
	fmt.Println(valsum)

	// Output:
	// 10
	// 30
}

func ExampleIFVector_Reduce() {
	vec := NewIfVector()
	vec[1] = 2
	vec[2] = 4
	vec[3] = 8

	sum := vec.Reduce(100, func(key uint64, curr, last float64) float64 {
		return last + curr
	})

	fmt.Println(sum)

	// Output:
	// 114
}

func ExampleIFVector_Map() {
	vec := NewIfVector()
	vec[10] = 0
	vec[20] = 0
	vec[30] = 0

	vec.Map(func(key uint64, val float64) float64 {
		return float64(key)
	})

	fmt.Println(vec[10])
	fmt.Println(vec[20])
	fmt.Println(vec[30])

	// Output:
	// 10
	// 20
	// 30
}

func ExampleIFVector_Filter() {
	vec := NewIfVector()
	vec[1] = 1
	vec[2] = 2
	vec[3] = 4
	vec[4] = 8

	vec.Filter(func(key uint64, val float64) bool {
		return key%2 == 0
	})

	// Kept because they were even.
	fmt.Println(vec[2])
	fmt.Println(vec[4])

	// Odd rows were removed
	_, ok := vec[1]
	fmt.Println(ok)

	_, ok = vec[3]
	fmt.Println(ok)

	// Output:
	// 2
	// 8
	// false
	// false
}

func ExampleIFVector_FilterMap() {
	vec := NewIfVector()
	vec[1] = 4
	vec[2] = 4

	filterOdd := func(key uint64, val float64) bool {
		return key%2 == 0
	}

	squareVal := func(key uint64, val float64) float64 {
		return val * val
	}

	vec.FilterMap(filterOdd, squareVal)

	fmt.Println(vec[1])
	fmt.Println(vec[2])

	// Output:
	// 0
	// 16
}

func ExampleIFVector_Sum() {
	vec := NewIfVector()
	vec[1] = 1
	vec[2] = 2
	vec[4] = 4

	sum := vec.Sum()

	fmt.Println(sum)

	// Output:
	// 7
}

func ExampleIFVector_Max() {
	vec := NewIfVector()
	vec[1] = 1
	vec[2] = 200
	vec[4] = -4

	max := vec.Max()

	fmt.Println(max)

	// Output:
	// 200
}

func ExampleIFVector_Dot() {
	a := NewIfVector()
	a[1] = -6
	a[2] = 8
	a[40] = 40 // no overlap

	b := NewIfVector()
	b[1] = 5
	b[2] = 12
	b[100] = 100

	fmt.Println(a.Dot(b))

	// Output:
	// 66
}

func BenchmarkIFVector_Dot1(bench *testing.B) {
	a := NewIfVector()
	b := NewIfVector()

	for i := 0; i < 500; i++ {
		if i%2 == 0 {
			b[uint64(i)] = float64(i)
		}
		a[uint64(i)] = float64(i)
	}

	for n := 0; n < bench.N; n++ {
		result64 = a.Dot(b)
	}
}

func BenchmarkIFVector_Dot2(bench *testing.B) {
	a := NewIfVector()
	b := NewIfVector()

	for i := 0; i < 500; i++ {
		if i%2 == 0 {
			b[uint64(i)] = float64(i)
		}
		a[uint64(i)] = float64(i)
	}

	for n := 0; n < bench.N; n++ {
		result64 = b.Dot(a)
	}
}

func BenchmarkIFVector_Dot3(bench *testing.B) {
	a := NewIfVector()
	b := NewIfVector()

	for i := 0; i < 500; i++ {
		b[uint64(i)] = float64(i)
		a[uint64(i)] = float64(i)
	}

	for n := 0; n < bench.N; n++ {
		result64 = b.Dot(a)
	}
}

func ExampleIFVector_L1Norm() {
	a := NewIfVector()
	a[1] = -4
	a[2] = 4

	fmt.Println(a.L1Norm())

	// Output:
	// 8
}

func ExampleIFVector_L2Norm() {
	a := NewIfVector()
	a[1] = -3
	a[2] = -2
	a[3] = 1

	fmt.Printf("%.3f\n", a.L2Norm())

	// Output:
	// 3.742
}

func ExampleIFVector_Clone() {
	a := NewIfVector()
	a[1] = 1
	a[2] = 2

	b := a.Clone()
	fmt.Println(b[1])
	fmt.Println(b[2])

	a[1] = 100
	fmt.Println(b[1])

	// Output:
	// 1
	// 2
	// 1
}

func ExampleIFVector_MultF() {
	a := NewIfVector()
	a[1] = 1
	a[2] = 2

	a.MultF(10)

	fmt.Println(a[1])
	fmt.Println(a[2])

	// Output:
	// 10
	// 20
}

func ExampleIFVector_DivF() {
	a := NewIfVector()
	a[1] = 1
	a[2] = 2

	a.DivF(2)

	fmt.Println(a[1])
	fmt.Println(a[2])

	// Output:
	// 0.5
	// 1
}

func ExampleIFVector_CosineSimilarity() {
	a := NewIfVector()
	a[1] = 0.702753576
	a[2] = 0.702753576

	b := NewIfVector()
	b[1] = 0.140550715
	b[2] = 0.140550715

	// do it manually
	dotProduct := a.Dot(b)
	l2A := a.L2Norm()
	l2B := b.L2Norm()
	manualResult := dotProduct / (l2A * l2B)
	fmt.Printf("%.3f\n", dotProduct)
	fmt.Printf("%.3f\n", l2A)
	fmt.Printf("%.3f\n", l2B)
	fmt.Printf("%.3f\n", manualResult)

	fmt.Printf("%.3f\n", a.CosineSimilarity(b))

	// Output:
	// 0.198
	// 0.994
	// 0.199
	// 1.000
	// 1.000
}

func TestOrderIFVectors(t *testing.T) {
	a := NewIfVector()
	a[0] = 1
	a[1] = 1
	a[2] = 1

	b := NewIfVector()
	b[0] = 1
	b[1] = 1

	var orderTests = []struct {
		first           IFVector
		second          IFVector
		expectedLarger  int
		expectedSmaller int
	}{
		{a, b, 3, 2},
		{b, a, 3, 2},
		{b, b, 2, 2},
		{a, a, 3, 3},
	}

	for _, test := range orderTests {
		smaller, larger := orderIFVectors(test.first, test.second)

		if len(smaller) != test.expectedSmaller {
			t.Errorf("len(%s): expected %d, actual %d", smaller, test.expectedSmaller, len(smaller))
		}

		if len(larger) != test.expectedLarger {
			t.Errorf("len(%s): expected %d, actual %d", larger, test.expectedLarger, len(larger))
		}
	}
}
