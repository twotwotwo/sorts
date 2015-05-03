// Copyright 2009 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package sortutil sorts and searches common slice types, and offers
// helper functions for sorting floats with radix sort.
package sortutil

import (
	"bytes"
	"github.com/twotwotwo/sorts"
	"math"
	"sort"
)

// Float32Key generates a uint64 key from a float32. To sort float32s,
// use this with Float32Less.
func Float32Key(f float32) uint64 {
	b := uint64(math.Float32bits(f)) << 32
	b ^= ^(b>>63 - 1) | (1 << 63)
	return b
}

// Float32Less compares float32s, sorting NaNs (which are normally
// unsortable) to the end.
func Float32Less(f, g float32) bool {
	return Float32Key(f) < Float32Key(g)
}

// Float64Key generates a uint64 key from a float64. To sort float64s,
// use this with Float64Less.
func Float64Key(f float64) uint64 {
	b := math.Float64bits(f)
	b ^= ^(b>>63 - 1) | (1 << 63)
	return b
}

// Float64Less compares float64s, sorting NaNs (which are normally
// unsortable) to the end.
func Float64Less(f, g float64) bool {
	return Float64Key(f) < Float64Key(g)
}

// IntSlice attaches the methods of Int64Interface to []int, sorting in increasing order.
type IntSlice []int

func (p IntSlice) Len() int           { return len(p) }
func (p IntSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p IntSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IntSlice) Key(i int) int64    { return int64(p[i]) }

// Sort is a convenience method.
func (p IntSlice) Sort() { sorts.ByInt64(p) }

// Int32Slice attaches the methods of Uint64Interface to []int32, sorting in increasing order.
type Int32Slice []int32

func (p Int32Slice) Len() int           { return len(p) }
func (p Int32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Int32Slice) Key(i int) int64    { return int64(p[i]) }

// Sort is a convenience method.
func (p Int32Slice) Sort() { sorts.ByInt64(p) }

// Int64Slice attaches the methods of Uint64Interface to []int64, sorting in increasing order.
type Int64Slice []int64

func (p Int64Slice) Len() int           { return len(p) }
func (p Int64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Int64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Int64Slice) Key(i int) int64    { return p[i] }

// Sort is a convenience method.
func (p Int64Slice) Sort() { sorts.ByInt64(p) }

// UintSlice attaches the methods of Uint64Interface to []uint, sorting in increasing order.
type UintSlice []uint

func (p UintSlice) Len() int           { return len(p) }
func (p UintSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p UintSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p UintSlice) Key(i int) uint64   { return uint64(p[i]) }

// Sort is a convenience method.
func (p UintSlice) Sort() { sorts.ByUint64(p) }

// Uint32Slice attaches the methods of Uint64Interface to []int32, sorting in increasing order.
type Uint32Slice []uint32

func (p Uint32Slice) Len() int           { return len(p) }
func (p Uint32Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Uint32Slice) Key(i int) uint64   { return uint64(p[i]) }

// Sort is a convenience method.
func (p Uint32Slice) Sort() { sorts.ByUint64(p) }

// Uint64Slice attaches the methods of Uint64Interface to []uint64, sorting in increasing order.
type Uint64Slice []uint64

func (p Uint64Slice) Len() int           { return len(p) }
func (p Uint64Slice) Less(i, j int) bool { return p[i] < p[j] }
func (p Uint64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Uint64Slice) Key(i int) uint64   { return p[i] }

// Sort is a convenience method.
func (p Uint64Slice) Sort() { sorts.ByUint64(p) }

// Float32Slice attaches the methods of Uint64Interface to []uint32, sorting in increasing order, NaNs last.
type Float32Slice []float32

func (p Float32Slice) Len() int           { return len(p) }
func (p Float32Slice) Less(i, j int) bool { return Float32Less(p[i], p[j]) }
func (p Float32Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float32Slice) Key(i int) uint64   { return Float32Key(p[i]) }

// Sort is a convenience method.
func (p Float32Slice) Sort() { sorts.ByUint64(p) }

// Float64Slice attaches the methods of Uint64Interface to []float64, sorting in increasing order, NaNs last.
type Float64Slice []float64

func (p Float64Slice) Len() int           { return len(p) }
func (p Float64Slice) Less(i, j int) bool { return Float64Less(p[i], p[j]) }
func (p Float64Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p Float64Slice) Key(i int) uint64   { return Float64Key(p[i]) }

// Sort is a convenience method.
func (p Float64Slice) Sort() { sorts.ByUint64(p) }

// StringSlice attaches the methods of StringInterface to []string, sorting in increasing order.
type StringSlice []string

func (p StringSlice) Len() int           { return len(p) }
func (p StringSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p StringSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p StringSlice) Key(i int) string   { return p[i] }

// Sort is a convenience method.
func (p StringSlice) Sort() { sorts.ByString(p) }

// BytesSlice attaches the methods of BytesInterface to [][]byte, sorting in increasing order.
type BytesSlice [][]byte

func (p BytesSlice) Len() int           { return len(p) }
func (p BytesSlice) Less(i, j int) bool { return bytes.Compare(p[i], p[j]) == -1 }
func (p BytesSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p BytesSlice) Key(i int) []byte   { return p[i] }

// Sort is a convenience method.
func (p BytesSlice) Sort() { sorts.ByBytes(p) }

// Ints sorts a slice of ints in increasing order.
func Ints(a []int) { IntSlice(a).Sort() }

// Int32s sorts a slice of int32s in increasing order.
func Int32s(a []int32) { Int32Slice(a).Sort() }

// Int64s sorts a slice of int64s in increasing order.
func Int64s(a []int64) { Int64Slice(a).Sort() }

// Uints sorts a slice of ints in increasing order.
func Uints(a []uint) { UintSlice(a).Sort() }

// Uint32s sorts a slice of uint64s in increasing order.
func Uint32s(a []uint32) { Uint32Slice(a).Sort() }

// Uint64s sorts a slice of uint64s in increasing order.
func Uint64s(a []uint64) { Uint64Slice(a).Sort() }

// Float32s sorts a slice of uint64s in increasing order, NaNs last.
func Float32s(a []float32) { Float32Slice(a).Sort() }

// Float64s sorts a slice of uint64s in increasing order, NaNs last.
func Float64s(a []float64) { Float64Slice(a).Sort() }

// Strings sorts a slice of strings in increasing order.
func Strings(a []string) { StringSlice(a).Sort() }

// Bytes sorts a slice of byte slices in increasing order.
func Bytes(a [][]byte) { BytesSlice(a).Sort() }

// IntsAreSorted tests whether a slice of ints is sorted in increasing order.
func IntsAreSorted(a []int) bool { return sort.IsSorted(IntSlice(a)) }

// Int32sAreSorted tests whether a slice of int32s is sorted in increasing order.
func Int32sAreSorted(a []int32) bool { return sort.IsSorted(Int32Slice(a)) }

// Int64sAreSorted tests whether a slice of int64s is sorted in increasing order.
func Int64sAreSorted(a []int64) bool { return sort.IsSorted(Int64Slice(a)) }

// UintsAreSorted tests whether a slice of ints is sorted in increasing order.
func UintsAreSorted(a []uint) bool { return sort.IsSorted(UintSlice(a)) }

// Uint32sAreSorted tests whether a slice of uint32s is sorted in increasing order.
func Uint32sAreSorted(a []uint32) bool { return sort.IsSorted(Uint32Slice(a)) }

// Uint64sAreSorted tests whether a slice of uint64s is sorted in increasing order.
func Uint64sAreSorted(a []uint64) bool { return sort.IsSorted(Uint64Slice(a)) }

// Float32sAreSorted tests whether a slice of float32s is sorted in increasing order, NaNs last.
func Float32sAreSorted(a []float32) bool { return sort.IsSorted(Float32Slice(a)) }

// Float64sAreSorted tests whether a slice of float64s is sorted in increasing order, NaNs last.
func Float64sAreSorted(a []float64) bool { return sort.IsSorted(Float64Slice(a)) }

// StringsAreSorted tests whether a slice of strings is sorted in increasing order.
func StringsAreSorted(a []string) bool { return sort.IsSorted(StringSlice(a)) }

// BytesAreSorted tests whether a slice of byte slices is sorted in increasing order.
func BytesAreSorted(a [][]byte) bool { return sort.IsSorted(BytesSlice(a)) }

// SearchInts searches ints; read about sort.Search for more.
func SearchInts(a []int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInts to the receiver and x.
func (p IntSlice) Search(x int) int { return SearchInts(p, x) }

// SearchInt32s searches int32s; read about sort.Search for more.
func SearchInt32s(a []int32, x int32) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInt32s to the receiver and x.
func (p Int32Slice) Search(x int32) int { return SearchInt32s(p, x) }

// SearchInt64s searches int64s; read about sort.Search for more.
func SearchInt64s(a []int64, x int64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchInt64s to the receiver and x.
func (p Int64Slice) Search(x int64) int { return SearchInt64s(p, x) }

// SearchUints searches uints; read about sort.Search for more.
func SearchUints(a []uint, x uint) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUints to the receiver and x.
func (p UintSlice) Search(x uint) int { return SearchUints(p, x) }

// SearchUint32s searches uint32s; read about sort.Search for more.
func SearchUint32s(a []uint32, x uint32) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUint32s to the receiver and x.
func (p Uint32Slice) Search(x uint32) int { return SearchUint32s(p, x) }

// SearchUint64s searches uint64s; read about sort.Search for more.
func SearchUint64s(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchUint64s to the receiver and x.
func (p Uint64Slice) Search(x uint64) int { return SearchUint64s(p, x) }

// SearchFloat32s searches float32s; read about sort.Search for more.
func SearchFloat32s(a []float32, x float32) int {
	return sort.Search(len(a), func(i int) bool { return Float32Key(a[i]) >= Float32Key(x) })
}

// Search returns the result of applying SearchFloat32s to the receiver and x.
func (p Float32Slice) Search(x float32) int { return SearchFloat32s(p, x) }

// SearchFloat64s searches float64s; read about sort.Search for more.
func SearchFloat64s(a []float64, x float64) int {
	return sort.Search(len(a), func(i int) bool { return Float64Key(a[i]) >= Float64Key(x) })
}

// Search returns the result of applying SearchFloat64s to the receiver and x.
func (p Float64Slice) Search(x float64) int { return SearchFloat64s(p, x) }

// SearchStrings searches strings; read about sort.Search for more.
func SearchStrings(a []string, x string) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// Search returns the result of applying SearchStrings to the receiver and x.
func (p StringSlice) Search(x string) int { return SearchStrings(p, x) }

// SearchBytes searches []bytes; read about sort.Search for more.
func SearchBytes(a [][]byte, x []byte) int {
	return sort.Search(len(a), func(i int) bool { return bytes.Compare(a[i], x) >= 0 })
}

// Search returns the result of applying SearchBytes to the receiver and x.
func (p BytesSlice) Search(x []byte) int { return SearchBytes(p, x) }
