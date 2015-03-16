package radixsort

import "sort"

// Copyright 2009 The Go Authors.
// Copyright 2014-5 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// This copies of code from sort.go because we can't use something like
// sort.SortRange(data, a, b) to sort a range of data.  Wrapping incoming
// data in another sort.Interface is possible, but kills speed.

// There's a small change to medianOfThree that reduces Swaps (though I
// didn't see a clear improvement on benchmarks from it) and exports an
// IsSorted that just calls stdlib sort's, for convenience.

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Insertion sort
func insertionSort(data sort.Interface, a, b int) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && data.Less(j, j-1); j-- {
			data.Swap(j, j-1)
		}
	}
}

// siftDown implements the heap property on data[lo, hi).
// first is an offset into the array where the root of the heap lies.
func siftDown(data sort.Interface, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data.Less(first+child, first+child+1) {
			child++
		}
		if !data.Less(first+root, first+child) {
			return
		}
		data.Swap(first+root, first+child)
		root = child
	}
}

func heapSort(data sort.Interface, a, b int) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data.Swap(first, first+i)
		siftDown(data, lo, i, first)
	}
}

// Quicksort, following Bentley and McIlroy,
// ``Engineering a Sort Function,'' SP&E November 1993.

// medianOfThree returns the middle of the three indicies
func medianOfThree(data sort.Interface, a, b, c int) (med int) {
	// If only one of a<b and a<c is true, a is the median. If only one
	// of a<c and b<c is true, c is the median. Otherwise, it's b.

	c0, c1 := data.Less(a, b), data.Less(a, c)
	// if c0 && !c1, then c <= a < b   (c <= a because !c1, a < b because c0)
	// if !c0 && c1, then b <= a < c   (b <= a because !c0, a < c because c1)
	if c0 != c1 {
		return a
	}

	c2 := data.Less(b, c)
	// if c1 && !c2, then a < c <= b   (a < c because c1, c <= b because !c2)
	// if !c1 && c2, then b < c <= a   (b < c because c2, c <= a because !c1)
	if c1 != c2 {
		return c
	}

	// if neither c0 != c1 or c1 != c2, c0 == c1 == c2.
	// and c0 == c2 leaves two possibilities:
	// if c0 && c2, then a < b < c        (a < b because c0, b < c because c2)
	// if !(c0 || c2), then c <= b <= a   (c <= b because !c2, b <= a because !c0)
	return b
}

func swapRange(data sort.Interface, a, b, n int) {
	for i := 0; i < n; i++ {
		data.Swap(a+i, b+i)
	}
}

func doPivot(data sort.Interface, lo, hi int) (midlo, midhi int) {
	m := lo + (hi-lo)/2 // Written like this to avoid integer overflow.
	m1, m2, m3 := lo, m, hi-1
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		m1 = medianOfThree(data, lo, lo+s, lo+2*s)
		m2 = medianOfThree(data, m, m-s, m+s)
		m3 = medianOfThree(data, hi-1, hi-1-s, hi-1-2*s)
	}
	data.Swap(lo, medianOfThree(data, m1, m2, m3))

	// Invariants are:
	//    data[lo] = pivot (set up by ChoosePivot)
	//    data[lo <= i < a] = pivot
	//    data[a <= i < b] < pivot
	//    data[b <= i < c] is unexamined
	//    data[c <= i < d] > pivot
	//    data[d <= i < hi] = pivot
	//
	// Once b meets c, can swap the "= pivot" sections
	// into the middle of the slice.
	pivot := lo
	a, b, c, d := lo+1, lo+1, hi, hi
	for {
		for b < c {
			if data.Less(b, pivot) { // data[b] < pivot
				b++
			} else if !data.Less(pivot, b) { // data[b] = pivot
				data.Swap(a, b)
				a++
				b++
			} else {
				break
			}
		}
		for b < c {
			if data.Less(pivot, c-1) { // data[c-1] > pivot
				c--
			} else if !data.Less(c-1, pivot) { // data[c-1] = pivot
				data.Swap(c-1, d-1)
				c--
				d--
			} else {
				break
			}
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] < pivot
		data.Swap(b, c-1)
		b++
		c--
	}

	n := min(b-a, a-lo)
	swapRange(data, lo, b-n, n)

	n = min(hi-d, d-c)
	swapRange(data, c, hi-n, n)

	return lo + b - a, hi - (d - c)
}

func quickSort(data sort.Interface, a, b, maxDepth int) {
	for b-a > 7 {
		if maxDepth == 0 {
			heapSort(data, a, b)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		insertionSort(data, a, b)
	}
}

// qSort quicksorts data.
// It performs O(n*log(n)) comparisons and swaps. The sort is not stable.
func qSort(data sort.Interface, a, b int) {
	// Switch to heapsort if depth of 2*ceil(lg(n+1)) is reached.
	n := b - a
	maxDepth := 0
	for i := n; i > 0; i >>= 1 {
		maxDepth++
	}
	maxDepth *= 2
	quickSort(data, a, b, maxDepth)
}

// IsSorted determines whether data is sorted.
func IsSorted(data sort.Interface) bool {
	return sort.IsSorted(data)
}
