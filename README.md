# in-mem-db In Memory Database
The project is intended to be easy to load and run. Save it to a folder on your computer, and navigate to its location in a Terminal.

I have provided both a zip file of the code (just in case, because I foolishly made the repo private) and this repo.

If:
```
$ git clone git@github.com:gwyndavidoff/in-mem-db.git
```
does not work then please download the zip file and use your preferred method to unzip.

# Requirements
Functioning Go: my environment had go 1.17

# Run the code
Both
```
$ go run src/main.go
$ make run
```
Will open up a prompt that accepts the following commands:
```
$ SET [name] [value] --requires at least 2 arguments, but ignores extra
$ GET [name] --requires at least 1 arguments, but ignores extra
$ DELETE [name] --requires at least 1 arguments, but ignores extra
$ COUNT [value] --requires at least 1 arguments, but ignores extra
$ END
$ BEGIN
$ ROLLBACK
$ COMMIT
```

You can also build before you run:
```
$ go build src/main.go
$ ./main
```

# Tests
There is a very small test suite that covers a couple of my initial sanity checks as well as the Examples from the spec. Both
```
$ go test ./...
$ make test
```
Will run and output the tests.

Some things I did not get to were:
* A more strenuous set of edge case tests
* I found a library that pretty prints btrees but I had some trouble getting the children from my nodes using the btree library and ran out of time [ppds](https://pkg.go.dev/github.com/shivamMg/ppds@v0.0.1/tree)

# Complexity and Requirements
I used the [btree] (https://pkg.go.dev/github.com/google/btree) library provided by google in golang. As a binary tree it should have runtime O log n for its operations, and I believe that my benchmarks indicate that the library fulfills that expectation.

In terms of the requirement not to double the db mem usage for every transaction, I wanted to note that btree actually does not fully copy the tree, and in fact is a lazy copy, but it seemed like that still wasn't in the spirit of what I understood the requirement to be about so I did not use it and instead kept lists of Rollback calls (also in a worst case if every entry in the db were part of the transaction it would eventually double the memory usage).

While I doubt the rollback call list is an ideal implementation (I'm pretty sure it's n log n or worse), I had considered doing something that would maintain the potential function calls during a transaction and then use those to make the changes on COMMIT and be available for searching during the transaction, but that would have slowed down retrieval while a transaction was open, so I opted for a more durable change to the saved data and a simple set of saved rollback commands. To improve it to at worst n log n I could change the set of commands to only contain at most one per key, though I'd need to double check that doesn't reduce performance of the SET and DELETE functions within a transaction. On the good side, I believe COMMIT is O(1).

Something else that was a bit of a tradeoff. I maintained a second tree of value counts. That allowed for the O log n access for the COUNT function, but in a worst case scenario would mean that the trees could be the same size.

I also do not check for equality, so "SET a foo" will perform a set of "a" to "foo", even if "a" is already "foo".

# Profiling
I added pprof on localhost:6060. Not sure how much there is to see, but if you go to [http://localhost:6060/debug/pprof/] (http://localhost:6060/debug/pprof/) in your browser you can explore it.

# Benchmark
Each Benchmark creates a db by adding N elements (where BenchmarkSetN), and then benchmarking an the matching operation (in these cases SET and GET). Out of curiosity I graphed both Set, Get, Count and Delete and they do appear to have the characteristics of a logarithmic graph.

I ran into some trouble with benchmarking Delete on a 0 size DB and have not yet debugged, which is why that value is missing.

To run the benchmarks you can call either of the below from the command line:
```
$ go test -bench=. ./...
$ make bench
```
Results
```
goos: darwin
goarch: amd64
pkg: in-mem-db/src/db
cpu: Intel(R) Core(TM) i9-9980HK CPU @ 2.40GHz
BenchmarkSet0-16           	 4527117	       245.2 ns/op
BenchmarkSet10-16          	 2436189	       488.9 ns/op
BenchmarkSet100-16         	 1212072	       997.6 ns/op
BenchmarkSet1000-16        	  760785	      1546 ns/op
BenchmarkSet10000-16       	  635214	      1901 ns/op
BenchmarkSet100000-16      	  443982	      2761 ns/op
BenchmarkSet1000000-16     	  378252	      3312 ns/op
BenchmarkGet0-16           	34507882	        36.69 ns/op
BenchmarkGet10-16          	13178354	        87.91 ns/op
BenchmarkGet100-16         	 7657562	       150.9 ns/op
BenchmarkGet1000-16        	 5633623	       201.0 ns/op
BenchmarkGet10000-16       	 4429179	       285.1 ns/op
BenchmarkGet100000-16      	 3487516	       341.1 ns/op
BenchmarkGet1000000-16     	 2685991	       444.8 ns/op
BenchmarkGet10000000-16    	 2425640	       503.1 ns/op
BenchmarkDelete10-16         	 1978672	       600.0 ns/op
BenchmarkDelete100-16        	  917733	      1357 ns/op
BenchmarkDelete1000-16       	  505618	      2515 ns/op
BenchmarkDelete10000-16      	  402015	      3108 ns/op
BenchmarkDelete100000-16     	  268780	      4190 ns/op
BenchmarkDelete1000000-16    	  240438	      5273 ns/op
BenchmarkCount0-16          	43665789	        28.39 ns/op
BenchmarkCount10-16         	13749255	        88.60 ns/op
BenchmarkCount100-16        	 7635684	       151.3 ns/op
BenchmarkCount1000-16       	 5959939	       224.0 ns/op
BenchmarkCount10000-16      	 4338177	       278.6 ns/op
BenchmarkCount100000-16     	 3358434	       341.7 ns/op
BenchmarkCount1000000-16    	 2787314	       447.3 ns/op
PASS
ok  	in-mem-db/src/db	174.635s
```