// Copyright 2013 The Go Authors.
// Copyright 2015 Randall Farmer.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package radixsort_test

import (
	"fmt"
	"github.com/twotwotwo/radixsort.test"
)

func Example_strings() {
	groceries := []string{"peppers", "tortillas", "tomatoes", "cheese"}
	radixsort.Strings(groceries) // or radixsort.Bytes([][]byte)
	fmt.Println(groceries)
	// Output: [cheese peppers tomatoes tortillas]
}
