// Copyright 2015 Randall Farmer. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radixsort_test

import . "github.com/twotwotwo/radixsort.test"
import "strconv"
import "sort"

// Helpers to sort in different ways: with the quicksort cutoff lowered or
// raised (to exercise radix code more or less), and sorting ints as bytes
// or strings (to exercise byte/string sort code).

// varyQSortCutoff runs a function with qSortCutoff set to 1, the default
// value, and 1e9.
func varyQSortCutoff(f func()) {
	f()
	defer SetQSortCutoff(SetQSortCutoff(1))
	f()
	SetQSortCutoff(1e9)
	f()
}

// forceRadix runs a function with qSortCutoff forced to 1
func forceRadix(f func()) {
	defer SetQSortCutoff(SetQSortCutoff(1))
	f()
}

// convertInts turns a []int of only positive numbers into a [][]byte and a
// []string ordered the same way, for exercising the string and []byte
// sorts.
func convertInts(a []int) ([][]byte, []string) { // and let's be honest it ain't no masterpiece either
	const l = 20 // length of converted number
	outSpace := make([]byte, l*len(a))
	for i := range outSpace {
		outSpace[i] = '0'
	}

	outBytes := make([][]byte, len(a))
	for i := range a {
		outBytes[i] = outSpace[l*i : l*i+l]
	}

	t := make([]byte, 20)
	for i, v := range a {
		s := strconv.AppendUint(t[:0], uint64(v), 10)
		copy(outBytes[i][l-len(s):], s)
	}

	strSpace := string(outSpace)
	outStrings := make([]string, len(a))
	for i := range a {
		outStrings[i] = strSpace[l*i : l*i+l]
	}

	return outBytes, outStrings
}

// multiSort sorts incoming integers using integer, []byte, and string sorts.
func multiSort(a []int) {
	asBytes, asStrings := convertInts(a)
	IntSlice(a).Sort()
	BytesSlice(asBytes).Sort()
	StringSlice(asStrings).Sort()
}

// manySort sorts integers with all QSort cutoffs and all data types. (It
// also throws away just tons of RAM but not sure we care.)
func manySort(a []int) {
	aBytes, aStrings := convertInts(a)

	myInts, myBytes, myStrings := IntSlice{}, BytesSlice{}, StringSlice{}
	varyQSortCutoff(func() {
		myInts = append(myInts[:0], a...)
		myInts.Sort()
		myBytes = append(myBytes[:0], aBytes...)
		myBytes.Sort()
		myStrings = append(myStrings[:0], aStrings...)
		myStrings.Sort()
	})
	// caller wants to see the sorted ints
	copy(a, myInts)
}

// manySortWrapper lets us use testBM with manySort to exercise all the sorts.
func manySortWrapper(d sort.Interface) {
	ints := d.(*testingData).data
	manySort(ints)
}
