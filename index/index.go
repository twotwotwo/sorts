// Copyright 2015 Randall Farmer. All rights reserved.

// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Package index uses sorted arrays of integers to assist sorting and
// searching, particularly of collections of strings.
package index

import (
	"sort"

	"github.com/twotwotwo/sorts"
)

type Index struct {
	Keys      []uint64
	Summary   []uint64 // implicit B-tree, if Summarize() was called
	Data      sort.Interface
	// intereted in these but not ready yet
	// Offset    int // index subslice where keys[0] is data[Offset]
	// KeyOffset int // for string keys that are 8th-15th bytes not 0-7th
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

// FindKey finds the position of the first item >= key in Keys, returning
// one after the end if there is none, and possibly using Summary for a
// speedup.  To find a range equal to k, you might use the range from
// FindKey(k) up to but not including FindKey(k+1) (with a special case
// using idx.Len() for k==^uint64(0), if relevant).  When different values
// map to the same key, you might want to sort.Search within the range to
// narrow your result down to the desired values.
func (idx *Index) FindKey(key uint64) int {
	if idx.Summary != nil {
		return idx.findKeySummary(key)
	}
	return sort.Search(idx.Len(), func(i int) bool { return idx.Keys[i] >= key })
}

// RF: awkwardness of docs above makes it mighty tempting to provide
// FindKeyRange and special cases for StringInterface and BytesInterface.

func (idx *Index) findKeySummary(key uint64) int {
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

// A type satisfying KeySetter can fill a []uint64 with numeric sort keys
// for its items.  Distinct values may map to the same integer key; Less
// is used as a tiebreaker.
type KeySetter interface {
	SetKeys([]uint64)
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
	case KeySetter:
		data.SetKeys(idx.Keys)
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
