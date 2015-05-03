// Copyright 2009 The Go Authors.
// Copyright 2014-5 Randall Farmer.
// All rights reserved.

// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radixsort_test

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"testing"

	. "github.com/twotwotwo/radixsort.test"
)

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var uints = [...]uint{74, 59, 238, 784, 9845, 959, 905, 0, 0, 42, 7586, 5467984, 7586}
var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, -1e30, 1e30, 905, 7.8, 7.8}
var strings = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestSortIntSlice(t *testing.T) {
	data := ints
	a := IntSlice(data[:])
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortInt32Slice(t *testing.T) {
	data := ints
	a := make(Int32Slice, len(data))
	for i, v := range data {
		a[i] = int32(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortInt64Slice(t *testing.T) {
	data := ints
	a := make(Int64Slice, len(data))
	for i, v := range data {
		a[i] = int64(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(-1e9) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUintSlice(t *testing.T) {
	data := uints
	a := UintSlice(data[:])
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", uints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUint32Slice(t *testing.T) {
	data := uints
	a := make(Uint32Slice, len(data))
	for i, v := range data {
		a[i] = uint32(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortUint64Slice(t *testing.T) {
	data := uints
	a := make(Uint64Slice, len(data))
	for i, v := range data {
		a[i] = uint64(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
	if a.Search(0) != 0 || a.Search(1e9) != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortFloat32Slice(t *testing.T) {
	data := float64s
	a := make(Float32Slice, len(data))
	for i, v := range data {
		a[i] = float32(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", a)
	}
	// IsSorted will compare using the Key, so compare using < to see if
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
	if a.Search(float32(math.Inf(-1))) != 0 || a.Search(float32(math.NaN())) != len(a)-2 {
		t.Errorf("search failed")
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := float64s
	a := Float64Slice(data[:])
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
	// IsSorted will compare using the Key, so compare using < to see if
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
	if a.Search(math.Inf(-1)) != 0 || a.Search(math.NaN()) != len(a)-2 {
		t.Errorf("search failed")
	}
}

func TestSortStringSlice(t *testing.T) {
	data := strings
	a := StringSlice(data[:])
	forceRadix(a.Sort)
	if !IsSorted(a) {
		t.Errorf("sorted %v", strings)
		t.Errorf("   got %v", data)
	}
	if a.Search("") != 0 || a.Search("\xFF") != len(a) {
		t.Errorf("search failed")
	}
}

func TestSortBytesSlice(t *testing.T) {
	dataStrings := strings
	a := make(BytesSlice, len(dataStrings))
	for i, v := range dataStrings {
		a[i] = []byte(v)
	}
	forceRadix(a.Sort)
	if !IsSorted(a) {
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

// check the IsSorted checks with a type that will never look sorted
type unsortableInts struct{ IntSlice }

func (u unsortableInts) Less(i, j int) bool { return j&1 == 1 }

type unsortableUints struct{ UintSlice }

func (u unsortableUints) Less(i, j int) bool { return j&1 == 1 }

type unsortableStrings struct{ StringSlice }

func (u unsortableStrings) Less(i, j int) bool { return j&1 == 1 }

type unsortableBytes struct{ BytesSlice }

func (u unsortableBytes) Less(i, j int) bool { return j&1 == 1 }

// more unsortable types, but now it's detectably because Key disagrees with Less
type miskeyedInts struct{ IntSlice }

func (u miskeyedInts) Less(i, j int) bool { return u.IntSlice[j] < u.IntSlice[i] }

type miskeyedUints struct{ UintSlice }

func (u miskeyedUints) Less(i, j int) bool { return u.UintSlice[j] < u.UintSlice[i] }

type miskeyedStrings struct{ StringSlice }

func (u miskeyedStrings) Less(i, j int) bool { return u.StringSlice[j] < u.StringSlice[i] }

type miskeyedBytes struct{ BytesSlice }

func (u miskeyedBytes) Less(i, j int) bool {
	return bytes.Compare(u.BytesSlice[j], u.BytesSlice[i]) == -1
}

func mustPanic(t *testing.T, name string, f func()) {
	defer func() { recover() }()
	f()
	t.Errorf("expected a panic on unsortable datatype %s", name)
}

func TestSortCheck(t *testing.T) {
	if !Checking() {
		return
	}
	mustPanic(t, "unsortableInts", func() {
		ByInt64(unsortableInts{IntSlice{1, 1, 1}})
	})
	mustPanic(t, "unsortableUints", func() {
		ByUint64(unsortableUints{UintSlice{1, 1, 1}})
	})
	mustPanic(t, "unsortableStrings", func() {
		ByString(unsortableStrings{StringSlice{"", "", ""}})
	})
	mustPanic(t, "unsortableBytes", func() {
		ByBytes(unsortableBytes{BytesSlice{[]byte{}, []byte{}, []byte{}}})
	})
	mustPanic(t, "miskeyedInts", func() {
		forceRadix(func() {
			ByInt64(miskeyedInts{IntSlice{1, 2, 3}})
		})
	})
	mustPanic(t, "miskeyedUints", func() {
		forceRadix(func() {
			ByUint64(miskeyedUints{UintSlice{1, 2, 3}})
		})
	})
	mustPanic(t, "miskeyedStrings", func() {
		forceRadix(func() {
			ByString(miskeyedStrings{StringSlice{"a", "b", "c"}})
		})
	})
	mustPanic(t, "miskeyedBytes", func() {
		forceRadix(func() {
			ByBytes(miskeyedBytes{BytesSlice{[]byte{'a'}, []byte{'b'}, []byte{'c'}}})
		})
	})
}

func TestFlip(t *testing.T) {
	data1, expected1 := [...]int{1, 2, 3, 4, 5}, [...]int{5, 4, 3, 2, 1}
	Flip(IntSlice(data1[:]))
	if data1 != expected1 {
		t.Errorf("Flip didn't flip!")
	}
	data2, expected2 := [...]int{1, 2}, [...]int{2, 1}
	Flip(IntSlice(data2[:]))
	if data2 != expected2 {
		t.Errorf("Flip didn't flip!")
	}
	Flip(IntSlice(nil)) // just shouldn't panic
}

func TestEmpty(t *testing.T) {
	IntSlice(nil).Sort()
	StringSlice(nil).Sort()
	BytesSlice(nil).Sort()
	IntSlice(nil).Search(0)
	StringSlice(nil).Search("")
	BytesSlice(nil).Search([]byte(nil))
}

func TestSortLarge_Random(t *testing.T) {
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	if IntsAreSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Ints(data)
	if !IntsAreSorted(data) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

func BenchmarkSortString1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]string, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = strconv.Itoa(i ^ 0x2cc)
		}
		b.StartTimer()
		Strings(data)
		b.StopTimer()
	}
}

func BenchmarkSortInt1K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<10)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0x2cc
		}
		b.StartTimer()
		Ints(data)
		b.StopTimer()
	}
}

func BenchmarkSortInt64K(b *testing.B) {
	b.StopTimer()
	for i := 0; i < b.N; i++ {
		data := make([]int, 1<<16)
		for i := 0; i < len(data); i++ {
			data[i] = i ^ 0xcccc
		}
		b.StartTimer()
		Ints(data)
		b.StopTimer()
	}
}

const (
	_Sawtooth = iota
	_Rand
	_Stagger
	_Plateau
	_Shuffle
	_NDist
)

const (
	_Copy = iota
	_Reverse
	_ReverseFirstHalf
	_ReverseSecondHalf
	_Sorted
	_Dither
	_NMode
)

type testingData struct {
	desc        string
	t           *testing.T
	data        []int
	maxswap     int // number of swaps allowed
	ncmp, nswap int
}

func (d *testingData) Len() int        { return len(d.data) }
func (d *testingData) Key(i int) int64 { return int64(d.data[i]) }
func (d *testingData) Less(i, j int) bool {
	d.ncmp++
	return d.data[i] < d.data[j]
}
func (d *testingData) Swap(i, j int) {
	if d.nswap >= d.maxswap {
		d.t.Errorf("%s: used %d swaps sorting slice of %d", d.desc, d.nswap, len(d.data))
		d.t.FailNow()
	}
	d.nswap++
	d.data[i], d.data[j] = d.data[j], d.data[i]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func lg(n int) int {
	i := 0
	for 1<<uint(i) < n {
		i++
	}
	return i
}

func testBentleyMcIlroy(t *testing.T, sort func(sort.Interface), maxswap func(int) int) {
	sizes := []int{100, 1023, 1024, 1025}
	if testing.Short() {
		sizes = []int{100, 127, 128, 129}
	}
	dists := []string{"sawtooth", "rand", "stagger", "plateau", "shuffle"}
	modes := []string{"copy", "reverse", "reverse1", "reverse2", "sort", "dither"}
	var tmp1, tmp2 [1025]int
	for _, n := range sizes {
		for m := 1; m < 2*n; m *= 2 {
			for dist := 0; dist < _NDist; dist++ {
				j := 0
				k := 1
				data := tmp1[0:n]
				for i := 0; i < n; i++ {
					switch dist {
					case _Sawtooth:
						data[i] = i % m
					case _Rand:
						data[i] = rand.Intn(m)
					case _Stagger:
						data[i] = (i*m + i) % n
					case _Plateau:
						data[i] = min(i, m)
					case _Shuffle:
						if rand.Intn(m) != 0 {
							j += 2
							data[i] = j
						} else {
							k += 2
							data[i] = k
						}
					}
				}

				mdata := tmp2[0:n]
				for mode := 0; mode < _NMode; mode++ {
					switch mode {
					case _Copy:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
					case _Reverse:
						for i := 0; i < n; i++ {
							mdata[i] = data[n-i-1]
						}
					case _ReverseFirstHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[n/2-i-1]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[i]
						}
					case _ReverseSecondHalf:
						for i := 0; i < n/2; i++ {
							mdata[i] = data[i]
						}
						for i := n / 2; i < n; i++ {
							mdata[i] = data[n-(i-n/2)-1]
						}
					case _Sorted:
						for i := 0; i < n; i++ {
							mdata[i] = data[i]
						}
						// Ints is known to be correct
						// because mode Sort runs after mode _Copy.
						Ints(mdata)
					case _Dither:
						for i := 0; i < n; i++ {
							mdata[i] = data[i] + i%5
						}
					}

					desc := fmt.Sprintf("n=%d m=%d dist=%s mode=%s", n, m, dists[dist], modes[mode])
					d := &testingData{desc: desc, t: t, data: mdata[0:n], maxswap: maxswap(n)}
					sort(d)
					// Uncomment if you are trying to improve the number of compares/swaps.
					//t.Logf("%s: ncmp=%d, nswp=%d", desc, d.ncmp, d.nswap)

					// If we were testing C qsort, we'd have to make a copy
					// of the slice and sort it ourselves and then compare
					// x against it, to ensure that qsort was only permuting
					// the data, not (for example) overwriting it with zeros.
					//
					// In go, we don't have to be so paranoid: since the only
					// mutating method Sort can call is TestingData.swap,
					// it suffices here just to check that the final slice is sorted.
					if !IntsAreSorted(mdata) {
						t.Errorf("%s: ints not sorted", desc)
						t.Errorf("\t%v", mdata)
						t.FailNow()
					}
				}
			}
		}
	}
}

func byInt64Wrapper(d sort.Interface) {
	ByInt64(d.(Int64Interface))
}

func TestSortBM(t *testing.T) {
	testBentleyMcIlroy(t, byInt64Wrapper, func(n int) int { return n * lg(n) * 12 / 10 })
}

func TestManySortBM(t *testing.T) {
	testBentleyMcIlroy(t, manySortWrapper, func(n int) int { return n * lg(n) * 12 / 10 })
}

func TestHeapsortBM(t *testing.T) {
	testBentleyMcIlroy(t, Heapsort, func(n int) int { return n * lg(n) * 12 / 10 })
}

// TestBackshift checks that radix sorting still works on data that trips up
// guessIntShift because it varies in a high bit, but only in a value that
// guessIntShift sampling misses
func TestBackshift(t *testing.T) {
	funnyData := [1e3]int{1: -1}
	funny := IntSlice(funnyData[:])
	if GuessIntShift(funny, len(funny)) > 0 {
		panic("guessIntShift got smarter")
	}
	forceRadix(func() { multiSort(funnyData[:]) })
	if !sort.IsSorted(funny) {
		t.Errorf("backshift data didn't sort")
	}
}

// TestFwdShift uses data that lets the radix sort shift past some bits in
// the middle; it might catch if it broke the sort.
func TestFwdShift(t *testing.T) {
	// an upper bit varies, lower byte varies, but bytes in between don't
	funnyData := []int{0x40000000, 23, 59, 38, 38, 6, 12, 9, 3, 4, 1, 49, 9, 63}
	funny := IntSlice(funnyData)
	forceRadix(func() { multiSort(funnyData) })
	if !sort.IsSorted(funny) {
		t.Errorf("forward-shift data didn't sort")
	}
}

// TestBrokenPrefix uses string and byte data where *most* input shares a
// common prefix except for one value that breaks the pattern at each byte
// position.  It's a bad case for the "everything was in one bucket"
// optimization, but we're merely looking for it not to barf (where barfing
// would be sort time exploding or data not sorting).
func TestBrokenPrefix(t *testing.T) {
	src := [128]byte{}
	src[64] = 1
	data := [10000][]byte{}
	for i := range data {
		data[i] = src[:]
	}
	// last 64 entries have a 1 in a pseudorandom position, breaking the
	// pattern
	for i := 10000 - 64; i < 10000; i++ {
		data[i] = src[64-((i*11)%64):]
	}
	forceRadix(BytesSlice(data[:]).Sort)
	if !BytesAreSorted(data[:]) {
		t.Errorf("broken-prefix data didn't sort")
	}

	srcStr := string(src[:])
	dataStr := [10000]string{}
	for i := range dataStr {
		dataStr[i] = srcStr
	}
	for i := 10000 - 64; i < 10000; i++ {
		data[i] = src[64-((i*11)%64):]
	}
	forceRadix(StringSlice(dataStr[:]).Sort)
	if !StringsAreSorted(dataStr[:]) {
		t.Errorf("broken-prefix data didn't sort")
	}
}

// TestShifts uses integer data consisting of a 1 bit in a random position.
// It's like TestBrokenPrefix for integer data.
func TestShifts(t *testing.T) {
	data := make([]uint64, 10000)
	for i := range data {
		data[i] = 1 << uint((i*19)%64)
	}
	forceRadix(Uint64Slice(data).Sort)
	if !Uint64sAreSorted(data) {
		t.Errorf("shifts data didn't sort")
	}
}

// TestMaxProcs makes sure forcing a serial sort doesn't break everything.
func TestMaxProcs(t *testing.T) {
	defer func(old int) { MaxProcs = old }(MaxProcs)
	MaxProcs = 1

	// this is TestLarge_Random
	n := 1000000
	if testing.Short() {
		n /= 100
	}
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	manySort(data)
	if !IntsAreSorted(data) {
		t.Errorf("serial sort failed")
	}
}

// TestSortByLength uses data that only varies in how many \0 bytes values
// contain.
func TestSortByLength(t *testing.T) {
	src := [128]byte{}
	data := [10000][]byte{}
	for i := range data {
		data[i] = src[:(i*19)%128]
	}
	forceRadix(BytesSlice(data[:]).Sort)
	if !BytesAreSorted(data[:]) {
		t.Errorf("sort-by-length data didn't sort")
	}

	srcStr := string(src[:])
	dataStr := [10000]string{}
	for i := range dataStr {
		dataStr[i] = srcStr[:(i*19)%128]
	}
	forceRadix(StringSlice(dataStr[:]).Sort)
	if !StringsAreSorted(dataStr[:]) {
		t.Errorf("sort-by-length data didn't sort")
	}
}

var countOpsSizes = []int{1e2, 3e2, 1e3, 3e3, 1e4, 3e4, 1e5, 3e5, 1e6}

func countOps(t *testing.T, algo func(sort.Interface), name string) {
	sizes := countOpsSizes
	if testing.Short() {
		sizes = sizes[:5]
	}
	if !testing.Verbose() {
		t.Skip("Counting skipped as non-verbose mode.")
	}
	for _, n := range sizes {
		td := testingData{
			desc:    name,
			t:       t,
			data:    make([]int, n),
			maxswap: 1<<31 - 1,
		}
		for i := 0; i < n; i++ {
			td.data[i] = rand.Intn(n / 5)
		}
		algo(&td)
		t.Logf("%s %8d elements: %11d Swap, %10d Less", name, n, td.nswap, td.ncmp)
	}
}

func TestCountSortOps(t *testing.T) { countOps(t, byInt64Wrapper, "Sort  ") }

func bench(b *testing.B, size int, algo func(sort.Interface), name string) {
	b.StopTimer()
	data := make(IntSlice, size)
	x := ^uint32(0)
	for i := 0; i < b.N; i++ {
		for n := size - 3; n <= size+3; n++ {
			for i := 0; i < len(data); i++ {
				x += x
				x ^= 1
				if int32(x) < 0 {
					x ^= 0x88888eef
				}
				data[i] = int(x % uint32(n/5))
			}
			b.StartTimer()
			algo(data)
			b.StopTimer()
			if !IsSorted(data) {
				b.Errorf("%s did not sort %d ints", name, n)
			}
		}
	}
}

func BenchmarkSort1e2(b *testing.B) { bench(b, 1e2, byInt64Wrapper, "Sort") }
func BenchmarkSort1e4(b *testing.B) { bench(b, 1e4, byInt64Wrapper, "Sort") }
func BenchmarkSort1e6(b *testing.B) { bench(b, 1e6, byInt64Wrapper, "Sort") }
