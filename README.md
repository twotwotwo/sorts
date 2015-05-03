# [radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test)

[radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test) provides
parallel radix sorting by a string, []byte, or (u)int64 key, and a parallel
quicksort.  Try it if faster sorting of large datasets is worth the hassle.

To radix sort, implement sort.Interface plus one more method, Key(i int),
returning the key for an item as string/[]byte/(u)int64, and call
radixsort.ByString, ByBytes, ByUint64, or ByInt64.  radixsort also exports
convenience functions and types like the stdlib's (Ints(data), IntSlice,
etc.), and a parallel Quicksort().  Change the value of radixsort.MaxProcs
to limit concurrency.  See the godoc for details and examples:
http://godoc.org/github.com/twotwotwo/radixsort.test

Note sort.Reverse() won't work here, but radixsort.Flip(data) will flip
ascending-sorted data to descending.  There's no stable sort.  The string
sorts just compare byte values; é and e won't sort next to each other.  Sort
functions that data is sorted after every run and panic(!) if it is not. 
You can sort floats with the functions like Float32Key and (to handle NaNs)
Float32Less.

Most folks should stick to stdlib sort: it's fast, standard, and easier to
use.  Still, radixsort can help if sorting huge datasets is a bottleneck for
you.

Love to hear if you're using this. You can reach me at my github username at
GMail, or on Twitter as @rf.
