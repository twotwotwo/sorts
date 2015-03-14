// Copyright 2010 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radixsort

import "sort"
import "bytes"

// Search calls sort.Search; consult its description.
func Search(n int, f func(int) bool) int { return sort.Search(n, f) }

// SearchInts searches ints; read about sort.Search for more.
func SearchInts(a []int, x int) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInts to the receiver and x.
func (p IntSlice) Search(x int) int { return SearchInts(p, x) }

// SearchInt32s searches int32s; read about sort.Search for more.
func SearchInt32s(a []int32, x int32) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInt32s to the receiver and x.
func (p Int32Slice) Search(x int32) int { return SearchInt32s(p, x) }

// SearchInt64s searches int64s; read about sort.Search for more.
func SearchInt64s(a []int64, x int64) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInt64s to the receiver and x.
func (p Int64Slice) Search(x int64) int { return SearchInt64s(p, x) }

// SearchUints searches uints; read about sort.Search for more.
func SearchUints(a []uint, x uint) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUints to the receiver and x.
func (p UintSlice) Search(x uint) int { return SearchUints(p, x) }

// SearchUint32s searches uint32s; read about sort.Search for more.
func SearchUint32s(a []uint32, x uint32) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUint32s to the receiver and x.
func (p Uint32Slice) Search(x uint32) int { return SearchUint32s(p, x) }

// SearchUint64s searches uint64s; read about sort.Search for more.
func SearchUint64s(a []uint64, x uint64) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUint64s to the receiver and x.
func (p Uint64Slice) Search(x uint64) int { return SearchUint64s(p, x) }

// SearchFloat32s searches float32s; read about sort.Search for more.
func SearchFloat32s(a []float32, x float32) int {
	return Search(len(a), func(i int) bool { return Float32Key(a[i]) >= Float32Key(x) })
}

// Search returns the result of applying SearchFloat32s to the receiver and x.
func (p Float32Slice) Search(x float32) int { return SearchFloat32s(p, x) }

// SearchFloat64s searches float64s; read about sort.Search for more.
func SearchFloat64s(a []float64, x float64) int {
	return Search(len(a), func(i int) bool { return Float64Key(a[i]) >= Float64Key(x) })
}

// Search returns the result of applying SearchFloat64s to the receiver and x.
func (p Float64Slice) Search(x float64) int { return SearchFloat64s(p, x) }

// SearchStrings searches strings; read about sort.Search for more.
func SearchStrings(a []string, x string) int {
	return Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchStrings to the receiver and x.
func (p StringSlice) Search(x string) int { return SearchStrings(p, x) }

// SearchBytes searches []bytes; read about sort.Search for more.
func SearchBytes(a [][]byte, x []byte) int {
	return Search(len(a), func(i int) bool { return bytes.Compare(a[i], x) >= 0 })
}

// Search returns the result of applying SearchBytes to the receiver and x.
func (p BytesSlice) Search(x []byte) int { return SearchBytes(p, x) }
