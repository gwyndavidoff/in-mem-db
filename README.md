# in-mem-db In Memory Database
The project is intended to be easy to load and run. Save it to a folder on your computer, and navigate to its location in a Terminal.

# Requirements
Functioning Go: my environment had go 1.17

# Run the code
Both
```
$ go run main.go
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

# Tests
There is a very small test suite that covers a couple of my initial sanity checks as well as the Examples from the spec. Both
```
$ go test ./...
$ make test
```
Will run and output the tests.

Some things I did not get to were:
* A more strenuous set of edge case tests
* Adding in pprof for profiling (may or may not have been helpful)
* Adding in benchmark tests. I discuss below why I believe I met the criteria of O log n for various methods, but I had intended to test and hopefully prove that with benchmarks

# Complexity and Requirements
I used the [btree] (https://pkg.go.dev/github.com/google/btree) library provided by google in golang. As a binary tree it should have runtime O log n for its operations, though with more time I had intended to use some stress/benchmark tests to verify whether this particular library fulfills that expectation.

In terms of the requirement not to double the db mem usage for every transaction, I wanted to note that btree actually does not fully copy the tree, and in fact is a lazy copy, but it seemed like that still wasn't in the spirit of what I understood the requirement to be about so I did not use it and instead kept lists of Rollback calls (also in a worst case if every entry in the db were part of the transaction it would eventually double the memory usage).

While I doubt the rollback call list is an ideal implementation (I'm pretty sure it's n log n or worse), I had considered doing something that would maintain the potential Function calls during a transaction and then use those to make the changes on COMMIT and be available for searching during the transaction, but that would have slowed down retrieval while a transaction was open, so I opted for a more durable change to the saved data and a simple set of saved rollback commands. To improve it to at worst n log n I could change the set of commands to only contain at most one per key.

Something else that was a bit of a tradeoff. I maintained a second tree of value counts. That allowed for the O log n access for the COUNT function, but in a worst case scenario would mean that the trees could be the same size.