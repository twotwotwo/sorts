// Copyright 2013 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sorts_test

import (
	"fmt"
	"github.com/twotwotwo/sorts"
	"github.com/twotwotwo/sorts/sortutil"
)

func Example_flip() {
	scores := []int{39, 492, 4912, 39, -10, 4, 92}
	data := sortutil.IntSlice(scores)
	data.Sort()
	sorts.Flip(data) // high scores first
	fmt.Println(scores)
	// Output: [4912 492 92 39 39 4 -10]
}
