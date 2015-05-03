// Copyright 2011 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radixsort

import "sort"

func Heapsort(data sort.Interface) {
	heapSort(data, 0, data.Len())
}

func GuessIntShift(data Int64Interface, l int) uint {
	return guessIntShift(intwrapper{data}, l)
}

func SetQSortCutoff(i int) int {
	orig := qSortCutoff
	qSortCutoff = i
	return orig
}

func SetMinOffload(i int) int {
	orig := minOffload
	minOffload = i
	return orig
}

func Checking() bool {
	return true
}
