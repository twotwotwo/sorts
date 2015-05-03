// Copyright 2013 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sorts_test

import (
	"fmt"
	"github.com/twotwotwo/sorts/sortutil"
)

func Example_strings() {
	groceries := []string{"peppers", "tortillas", "tomatoes", "cheese"}
	sortutil.Strings(groceries) // or sortutil.Bytes([][]byte)
	fmt.Println(groceries)
	// Output: [cheese peppers tomatoes tortillas]
}
