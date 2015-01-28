package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

// Note - the sorting algos might mutate input slices. The func([]int)[]int sig
// is just to keep things consistent.

// Number of comparisons.
const (
	n       = 103
	repeats = 200
)

var (
	countCompares = true
	compareCount  = 0
)

var (
	//sortAlgo      = qsort
	//sortAlgo      = msort
	sortAlgo      = goSort
	//sortAlgo      = isort
)

func init() {
	// Defaults to seed == 1 if you don't specify...!!
	rand.Seed(time.Now().UnixNano())
}

// Return -ve for n1 < n2, 0 for n1 == n2, +ve for n1 > n2
func compare(n1, n2 int) int {
	if countCompares {
		compareCount++
	}

	return n1 - n2
}

// Go Sort --------------------

type sortableInts []int

func (ns sortableInts) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns sortableInts) Len() int {
	return len(ns)
}

func (ns sortableInts) Less(i, j int) bool {
	return compare(ns[i], ns[j]) <= 0
}

func goSort(ns []int) []int {
	sort.Sort(sortableInts(ns))
	return ns
}

// Merge Sort --------------------

func merge(ns1, ns2 []int) []int {
	len1, len2 := len(ns1), len(ns2)
	total := len1 + len2
	i, j := 0, 0

	ret := make([]int, total)

	for k := 0; k < total; k++ {
		// One of the two inputs are exhausted.
		if i >= len1 {
			copy(ret[k:], ns2[j:])
			break
		} else if j >= len2 {
			copy(ret[k:], ns1[i:])
			break
		}

		diff := compare(ns1[i], ns2[j])
		if diff <= 0 {
			ret[k] = ns1[i]
			i++
		} else {
			ret[k] = ns2[j]
			j++
		}
	}

	return ret
}

func msort(ns []int) []int {
	if len(ns) <= 1 {
		return ns
	}

	mid := len(ns) / 2

	return merge(msort(ns[:mid]), msort(ns[mid:]))
}

// Insertion Sort --------------------

func isort(ns []int) []int {
	for i := 1; i < len(ns); i++ {
		key := ns[i]
		j := i - 1
		for ; j >= 0 && compare(ns[j], key) > 0; j-- {
			ns[j+1] = ns[j]
		}
		ns[j+1] = key
	}

	return ns
}

// Quick Shit Sort --------------------

func getPivot(ns []int, from, to int) int {
	// Calculate midpoint without overflow. This is a shit pivot.
	return to - (to-from)/2
}

func swap(ns []int, i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func partition(ns []int, from, to, pivot int) int {
	pivotVal := ns[pivot]

	// Pop the pivot at the end of the array...
	swap(ns, pivot, to)

	// ...Reset the pivot to the start...
	pivot = from

	for i := from; i < to; i++ {
		if compare(ns[i], pivotVal) < 0 {
			swap(ns, i, pivot)
			pivot++
		}
	}

	// ...And pop it back again!
	swap(ns, to, pivot)

	return pivot
}

func qsort(ns []int) []int {
	var qs func(from, to int)

	qs = func(from, to int) {
		if from >= to {
			return
		}

		pivot := getPivot(ns, from, to)
		pivot = partition(ns, from, to, pivot)
		qs(from, pivot-1)
		qs(pivot+1, to)
	}

	qs(0, len(ns)-1)

	return ns
}

// Testing Code --------------------

func genRange() []int {
	ret := make([]int, n)

	for i := 0; i < n; i++ {
		ret[i] = int(rand.Int31n(n))
	}

	return ret
}

func assertSorted(ns []int) {
	currCountCompares := countCompares
	countCompares = false
	defer func() {
		countCompares = currCountCompares
	}()

	if len(ns) <= 1 {
		return
	}

	for i := 0; i < len(ns)-1; i++ {
		if compare(ns[i], ns[i+1]) > 0 {
			panic("not sorted!")
		}
	}
}

func run() int {
	defer func() {
		compareCount = 0
	}()

	// Generate range of (probably :) unsorted numbers and sort.
	ns := genRange()
	assertSorted(sortAlgo(ns))

	return compareCount
}

func estimate() float64 {
	return float64(n) * math.Log2(float64(n))
}

func main() {
	fmt.Printf("~ %0.0f\n", estimate())

	min, max := int(^uint(0)>>1), -1

	av := 0

	for i := 0; i < repeats; i++ {
		count := run()

		if count < min {
			min = count
		}
		if count > max {
			max = count
		}

		// Ref: http://math.stackexchange.com/a/106314
		av = (i*av + count) / (i + 1)

		if i%10 == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("%d\t", count)
	}
	fmt.Printf("\n")

	fmt.Printf("\n%d\t%d\t%d\n", min, av, max)
}
