// Copyright 2015 Randall Farmer. All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package index uses sorted arrays of integers to assist sorting and
// searching, particularly of collections of strings.
package index

import (
	"bytes"
	"sort"
	"strings"

	"github.com/twotwotwo/sorts"
)

type Index struct {
	Keys    []uint64
	Summary []uint64 // implicit B-tree, if Summarize() was called
	Data    sort.Interface
}

// Len returns the length of the data underlying an Index
func (idx *Index) Len() int { return idx.Data.Len() }

// Less compares Index elements by their Keys, falling back to Data.Less for
// equal-keyed items.
func (idx *Index) Less(i, j int) bool {
	return idx.Keys[i] < idx.Keys[j] || (idx.Keys[i] == idx.Keys[j] && idx.Data.Less(i, j))
}

// Swap swaps both the Keys and the inderlying data items at indices i and
// j.
func (idx *Index) Swap(i, j int) {
	idx.Keys[i], idx.Keys[j] = idx.Keys[j], idx.Keys[i]
	idx.Data.Swap(i, j)
}

// Key returns the uint64 key at index i.
func (idx *Index) Key(i int) uint64 { return idx.Keys[i] }

// levelBits and pageSize control the fan-out of Summary, the implicit
// B-tree.  6 won a very informal bake-off.  (Would have guessed 3, matching
// 8-word amd64 cache lines.) More would work better if this were ever on
// block storage, e.g.  levelBits of 9 corresponds to a page size of 4KiB.
const levelBits = 6
const pageSize = 1 << levelBits

// Summarize makes an implicit B-tree to speed lookups, using a few percent
// overhead on top of what's already used for Indices.
func (idx *Index) Summarize() {
	l := idx.Len()
	sl := l>>levelBits + l>>levelBits*2 + l>>levelBits*3 + l>>((levelBits*4)-1)
	summary := make([]uint64, 0, sl)
	summarizing := idx.Keys
	levelNum := 1
	for len(summarizing) > pageSize {
		start := len(summary)
		for i := 0; i < len(summarizing); i += pageSize {
			summary = append(summary, summarizing[i])
		}
		summarizing = summary[start:]
		levelNum++
	}
	idx.Summary = summary
}

// FindUint64 finds the position of the first item >= key in Keys, returning
// one after the end if there is none.  When different values map to the same key,
// you might want to sort.Search within the returned range to narrow your result
// down to the desired values.
func (idx *Index) FindUint64(key uint64) int {
	if idx.Summary != nil {
		return idx.findUint64Summary(key)
	}
	return sort.Search(idx.Len(), func(i int) bool { return idx.Keys[i] >= key })
}

// Compares string a to []byte b, returning -1 if a<b, 0 if a==b, and 1 if a>b.
func CompareStringToBytes(a string, b []byte) int {
	for i := range b {
		if i > len(a) {
			return -1
		}
		if b[i] > a[i] {
			return -1
		}
		if a[i] > b[i] {
			return 1
		}
	}
	if len(a) > len(b) {
		return 1
	}
	return 0
}

// CompareBytesToString is a convenience wrapper for CompareStringToBytes.
func CompareBytesToString(a []byte, b string) int {
	return -CompareStringToBytes(b, a)
}

// FindString finds the first item >= key, returning one after the end if there
// is none. The collection type must implement Key(i) returning string or []byte.
func (idx *Index) FindString(key string) int {
	k := StringKey(key)
	a, b := idx.FindUint64Range(k)
	switch data := idx.Data.(type) {
	case sorts.StringInterface:
		return a + sort.Search(b-a, func(i int) bool {
			return strings.Compare(key, data.Key(a+i)) >= 0
		})
	case sorts.BytesInterface:
		return a + sort.Search(b-a, func(i int) bool {
			return CompareStringToBytes(key, data.Key(a+i)) >= 0
		})
	default:
		panic("to use FindStringKey, Data.Key(i) must return string or []byte")
	}
}

// FindBytes finds the first item >= key, returning one after the end if there
// is none. The collection type must implement Key(i) returning string or []byte.
func (idx *Index) FindBytes(key []byte) int {
	k := BytesKey(key)
	a, b := idx.FindUint64Range(k)
	switch data := idx.Data.(type) {
	case sorts.StringInterface:
		return a + sort.Search(b-a, func(i int) bool {
			return CompareBytesToString(key, data.Key(a+i)) >= 0
		})
	case sorts.BytesInterface:
		offset := sort.Search(b-a, func(i int) bool {
			return bytes.Compare(key, data.Key(a+i)) >= 0
		})
		return a + offset
	default:
		panic("to use FindBytesKey, Data.Key(i) must return string or []byte")
	}
}

// FindUint64Range looks for a range of keys such that all items in idx.Keys[a:b]
// equal key.
// It can return an empty range if the item isn't found; in that case, a and b are both where the item would be inserted (and can be one past the end).
// To find a single key, use FindUint64.
func (idx *Index) FindUint64Range(key uint64) (a, b int) {
	a = idx.FindUint64(key)
	if a == len(idx.Keys) || idx.Keys[a] != key {
		// key not found, needn't search again
		return a, a
	}
	if key == ^uint64(0) { // would overflow
		b = len(idx.Keys)
	} else {
		b = idx.FindUint64(key + 1)
	}
	return
}

// FindStringRange(key) finds the range (a,b] such that Key() returns key for all items in idx.Data[a:b].
// It can return an empty range if the item isn't found; in that case, a and b are both where the item would be inserted (and can be one past the end).
// Data must implement Key(i) returning string or []byte.
// To find a single item, use FindString.
func (idx *Index) FindStringRange(key string) (int, int) {
	k := StringKey(key)
	a, b := idx.FindUint64Range(k)
	switch data := idx.Data.(type) {
	case sorts.StringInterface:
		aa := a + sort.Search(b-a, func(i int) bool {
			return strings.Compare(key, data.Key(a+i)) >= 0
		})
		bb := aa + sort.Search(b-aa, func(i int) bool {
			return strings.Compare(key, data.Key(aa+i)) > 0
		})
		return aa, bb
	case sorts.BytesInterface:
		aa := a + sort.Search(b-a, func(i int) bool {
			return CompareStringToBytes(key, data.Key(a+i)) >= 0
		})
		bb := aa + sort.Search(b-aa, func(i int) bool {
			return CompareStringToBytes(key, data.Key(aa+i)) > 0
		})
		return aa, bb
	default:
		panic("to use FindStringRange, Data.Key(i) must return string or []byte")
	}
}

// FindBytesRange(key) finds the range (a,b] such that Key() returns key for all items in idx.Data[a:b].
// It can return an empty range if the item isn't found; in that case, a is where the item would be inserted (and can be one past the end).
// Data must implement Key(i) returning string or []byte.
// To find a single item, use FindBytes.
func (idx *Index) FindBytesRange(key []byte) (int, int) {
	k := BytesKey(key)
	a, b := idx.FindUint64Range(k)
	switch data := idx.Data.(type) {
	case sorts.StringInterface:
		aa := a + sort.Search(b-a, func(i int) bool {
			return CompareBytesToString(key, data.Key(a+i)) >= 0
		})
		bb := aa + sort.Search(b-aa, func(i int) bool {
			return CompareBytesToString(key, data.Key(aa+i)) > 0
		})
		return aa, bb
	case sorts.BytesInterface:
		aa := a + sort.Search(b-a, func(i int) bool {
			return bytes.Compare(key, data.Key(a+i)) >= 0
		})
		bb := aa + sort.Search(b-aa, func(i int) bool {
			return bytes.Compare(key, data.Key(aa+i)) > 0
		})
		return aa, bb
	default:
		panic("to use FindStringRangeExact, Data.Key(i) must return string or []byte")
	}
}

func (idx *Index) findUint64Summary(key uint64) int {
	summary := idx.Summary
	keys := idx.Keys

	// count how many layers to expect in the "btree"
	levels, l := 0, len(keys)
	for l > 0 {
		levels++
		l >>= levelBits
	}
	levels-- // went one too far

	// keep following largest-strictly-less down the chain
	levelNum := levels
	levelEnd := len(summary)
	offset := 0
	for levelNum > 0 {
		// extract the "level"
		thisLevelBits := uint(levelBits * levelNum)
		levelLen := len(keys) >> thisLevelBits
		if len(keys) > levelLen<<thisLevelBits {
			// an entry for the remainder
			levelLen++
		}
		level := summary[levelEnd-levelLen : levelEnd]

		// extract the page at the given offset
		pageEnd := offset + pageSize
		if pageEnd > len(level) {
			pageEnd = len(level)
		}
		page := level[offset:pageEnd]

		// scan page for an entry >= key
		// binsearch would be fewer operations but less predictable ones
		i := 0
		for i < len(page) && page[i] < key {
			i++
		}
		if i > 0 {
			// i is first one that goes too far
			i--
		}

		// use that to walk down the tree
		// in particular, get next offset and level location
		offset += i
		offset <<= levelBits
		levelEnd -= levelLen
		levelNum--
	}

	// level==0 is the original array
	pageEnd := offset + pageSize
	if pageEnd > len(keys) {
		pageEnd = len(keys)
	}
	page := keys[offset:pageEnd]
	i := 0
	for i < len(page) && page[i] < key {
		i++
	}
	return offset + i
}

// StringKey generates a uint64 key from the first bytes of key.
func StringKey(key string) uint64 {
	k := uint64(0)
	for j := 0; j < 8 && j < len(key); j++ {
		k ^= uint64(key[j]) << uint(56-8*j)
	}
	return k
}

// BytesKey generates a uint64 key from the first bytes of key.
func BytesKey(key []byte) uint64 {
	k := uint64(0)
	for j := 0; j < 8 && j < len(key); j++ {
		k ^= uint64(key[j]) << uint(56-8*j)
	}
	return k
}

// SortWithIndex allocates an Index with space for a uint64 key for each
// item in data, then sorts items by their uint64 keys, using data.Less as a
// tie-breaker for equal-keyed items.  data may implement index.KeySetter or
// any of sorts.StringInterface, BytesInterface, or Uint64Interface.
func SortWithIndex(data sort.Interface) *Index {
	l := data.Len()
	indices := make([]uint64, l)
	idx := &Index{
		Keys: indices,
		Data: data,
	}
	switch data := data.(type) {
	case sorts.StringInterface:
		for i := 0; i < l; i++ {
			key := data.Key(i)
			k := uint64(0)
			for j := 0; j < 8 && j < len(key); j++ {
				k ^= uint64(key[j]) << uint(56-8*j)
			}
			indices[i] = k
		}
	case sorts.BytesInterface:
		for i := 0; i < l; i++ {
			key := data.Key(i)
			k := uint64(0)
			for j := 0; j < 8 && j < len(key); j++ {
				k ^= uint64(key[j]) << uint(56-8*j)
			}
			indices[i] = k
		}
	case sorts.Uint64Interface:
		for i := 0; i < l; i++ {
			indices[i] = data.Key(i)
		}
	default:
		panic("don't know how to extract int keys for data")
	}
	sorts.ByUint64(idx)
	return idx
}
