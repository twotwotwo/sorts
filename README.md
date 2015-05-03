# [radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test)

[radixsort](http://godoc.org/github.com/twotwotwo/radixsort.test) provides
parallel radix sorting by a string, []byte, or (u)int64 key, and a parallel
Quicksort(data). 
[sortutil](http://godoc.org/github.com/twotwotwo/radixsort.test/sortutil)
sorts common slice types and adds functions to help sort floats.

Usually, stick to stdlib sort: that's fast, standard, and simpler.  But
radixsort can help if sorting huge datasets is a bottleneck for you.  Try it
if shorter sort times seem worth some hassle in your application.

To radix sort, 
[implement sort.Interface](http://golang.org/pkg/sort/#Interface) 
plus one more method, Key(i int), returning the key for an item as
string/[]byte/(u)int64, and call radixsort.ByString, ByBytes, ByUint64, or
ByInt64.  Set radixsort.MaxProcs if you want to limit concurrency.  See the
godoc for details and examples:
http://godoc.org/github.com/twotwotwo/radixsort.test

There's no Reverse(), but radixsort.Flip(data) will flip ascending-sorted
data to descending.  There's no stable sort.  The string sorts just compare
byte values; so é won't sort next to e.  The package checks that data is
sorted after every run and panics(!) if not.

I'd love to hear if you're using this. E-mail me at my github username at
GMail, or contact me on Twitter ([@rf](http://twitter.com/rf/)).
