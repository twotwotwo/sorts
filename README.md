# [radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test)

[radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test) provides radix sorting by a string, []byte, or uint64 key. It has helpers for sorting other numeric types (floats and signed ints) with the uint64 sort. You'd use it as a faster alternative to the standard package sort in places where you're sorting at least thousands of items by a key of one of the supported types, and the gain in sort speed matters enough to justify using a less elegant interface than the stdlib's.

Use it by implementing sort.Interface plus one more method, Key(i int), returning the key for an item as string/[]byte/uint64. Then you call radixsort.ByString, ByBytes, or ByNumber on your data as appropriate. The package exports convenience types and functions like the stdlib's does (radixsort.Ints(), etc.). See the godoc for details and examples: http://godoc.org/github.com/twotwotwo/radixsort.test

Some usage gotchas: when sorting float or signed int data, remember to use the helper functions like Int32Key or Float32Key and Float32Less; just converting numbers to uint64 will mishandle negative numbers, and not-a-number float values are tricky. sort.Reverse() won't work with radixsort, but radixsort.Flip(data) will flip ascending-sorted data to descending. There's no stable sort. Remember that the string sorts aren't doing any fancy collation, just handling data as raw bytes, so pairs like é and e won't sort next to each other.

Most folks just should use stdlib sort: there's a good chance sorting is not a bottleneck for you, and thus not worth extra effort to speed up; also, radixsort only helps when sorting at least thousands of elements at once, rather than doing lots of smaller sorts. Speedups vary with how your data looks--lots of duplicate items, or very little entropy in the early bytes of strings, slow things down a bit.

-----

The ".test" reflects that this may not be the eventual API I'll put up as plain "radixsort", and also that this is a newly released library and, though it has tests, I'd love to get more confidence in its correctness from folks trying it out. Right now it always checks that your results are really sorted and panics(!) if not. If you want to report a problem with radixsort, include your Less and Key functions, since mistakes in them (or data races) can be hard to distinguish from bugs in radixsort.

Love to hear if you're using this. You can reach me at my github username at GMail, or on Twitter as @rf.
