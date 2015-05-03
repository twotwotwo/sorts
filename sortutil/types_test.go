// Copyright 2009 The Go Authors.
// Copyright 2014-5 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sortutil_test

import (
	. "github.com/twotwotwo/sorts/sortutil"
	"math"
	"sort"
	"testing"
)

// we need enough elements that radix sort will kick in, or we're not
// really testing our Key implementations at all.
var testSize = 1024

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var uints = [...]uint{74, 59, 238, 784, 9845, 959, 905, 0, 0, 42, 7586, 5467984, 7586}
var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, -1e30, 1e30, 905, 7.8, 7.8}
var strings = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestSortIntSlice(t *testing.T) {
	data := ints
	a := make(IntSlice, testSize)
	for i := range a {
		a[i] = data[i%len(data)]
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortInt32Slice(t *testing.T) {
	data := ints
	a := make(Int32Slice, testSize)
	for i := range a {
		a[i] = int32(data[i%len(data)])
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortInt64Slice(t *testing.T) {
	data := ints
	a := make(Int64Slice, testSize)
	for i := range a {
		a[i] = int64(data[i%len(data)])
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUintSlice(t *testing.T) {
	data := uints
	a := make(UintSlice, testSize)
	for i := range a {
		a[i] = uint(data[i%len(data)])
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", uints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUint32Slice(t *testing.T) {
	data := uints
	a := make(Uint32Slice, testSize)
	for i := range a {
		a[i] = uint32(data[i%len(data)])
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUint64Slice(t *testing.T) {
	data := uints
	a := make(Uint64Slice, testSize)
	for i := range a {
		a[i] = uint64(data[i%len(data)])
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortFloat32Slice(t *testing.T) {
	data := float64s
	a := make(Float32Slice, testSize)
	nanCount := 0
	for i := range a {
		a[i] = float32(data[i%len(data)])
		if math.IsNaN(float64(a[i])) {
			nanCount++
		}
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", a)
	}
	// sort.IsSorted will compare using the Key, so compare using < to see if
	// Key is wrong
	prev := a[0]
	for _, v := range a {
		if v < prev {
			t.Errorf("Float32Key is wrong: %f sorted before %f", prev, v)
		}
		prev = v
	}
	// floats data contains two NaNs, so Search will find the spot right
	// before them.
	if a.Search(float32(math.Inf(-1))) != 0 || a.Search(float32(math.NaN())) != len(a)-nanCount {
		t.Errorf("search failed")
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := float64s
	a := make(Float64Slice, testSize)
	nanCount := 0
	for i := range a {
		a[i] = data[i%len(data)]
		if math.IsNaN(a[i]) {
			nanCount++
		}
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
	// sort.IsSorted will compare using the Key, so compare using < to see if
	// Key func is wrong
	prev := a[0]
	for _, v := range a {
		if v < prev {
			t.Errorf("Float64Key is wrong: %f sorted before %f", prev, v)
		}
		prev = v
	}
	// floats data contains two NaNs, so Search will find the spot right
	// before them.
	if a.Search(math.Inf(-1)) != 0 || a.Search(math.NaN()) != len(a)-nanCount {
		t.Errorf("search failed")
	}
}

func TestSortStringSlice(t *testing.T) {
	data := strings
	a := make(StringSlice, testSize)
	for i := range a {
		a[i] = data[i%len(data)]
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
	if a.Search("") != 0 || a.Search("\xFF") != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortBytesSlice(t *testing.T) {
	dataStrings := strings
	data := [][]byte{}
	for _, v := range dataStrings {
		data = append(data, []byte(v))
	}
	a := make(BytesSlice, testSize)
	for i := range a {
		a[i] = data[i%len(data)]
	}
	a.Sort()
	if !sort.IsSorted(a) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", a)
	}
	if a.Search([]byte(nil)) != 0 || a.Search([]byte{255}) != len(a) {
		t.Errorf("search failed")
	}
}

func TestInts(t *testing.T) {
	data := ints
	Ints(data[:])
	if !IntsAreSorted(data[:]) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestInt32s(t *testing.T) {
	data := make([]int32, len(ints))
	for i, v := range ints {
		data[i] = int32(v)
	}
	Int32s(data)
	if !Int32sAreSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestInt64s(t *testing.T) {
	data := make([]int64, len(ints))
	for i, v := range ints {
		data[i] = int64(v)
	}
	Int64s(data)
	if !Int64sAreSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestUints(t *testing.T) {
	data := uints
	Uints(data[:])
	if !UintsAreSorted(data[:]) {
		t.Errorf("sorted %v", uints)
		t.Errorf("   got %v", data)
	}
}

func TestUint32s(t *testing.T) {
	data := make([]uint32, len(uints))
	for i, v := range uints {
		data[i] = uint32(v)
	}
	Uint32s(data)
	if !Uint32sAreSorted(data) {
		t.Errorf("sorted %v", uints)
		t.Errorf("   got %v", data)
	}
}

func TestUint64s(t *testing.T) {
	data := make([]uint64, len(uints))
	for i, v := range uints {
		data[i] = uint64(v)
	}
	Uint64s(data)
	if !Uint64sAreSorted(data) {
		t.Errorf("sorted %v", uints)
		t.Errorf("   got %v", data)
	}
}

func TestFloat32s(t *testing.T) {
	data := make([]float32, len(float64s))
	for i, v := range float64s {
		data[i] = float32(v)
	}
	Float32s(data)
	if !Float32sAreSorted(data) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestFloat64s(t *testing.T) {
	data := make([]float64, len(float64s))
	for i, v := range float64s {
		data[i] = v
	}
	Float64s(data)
	if !Float64sAreSorted(data) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestStrings(t *testing.T) {
	data := strings
	Strings(data[:])
	if !StringsAreSorted(data[:]) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
}

func TestBytes(t *testing.T) {
	data := make([][]byte, len(strings))
	for i, v := range strings {
		data[i] = []byte(v)
	}
	Bytes(data[:])
	if !BytesAreSorted(data[:]) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
}
