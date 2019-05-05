# go-kmv

**go-kmv** is an adaptive version of *K-minimum values algorithm for cardinality estimation*

This repository provides:
  - A **library** for your own Go programs
  - A **cmd tool** which estimates the cardinality reading from the stdin (so you can use it with the pipe `|` linux operator)

# Examples

After compiling `cmd/main.go` we can run the algorithm from our terminal

```bash
$ go build -o go-kmv main.go

# Output
# ${CardinalityEstimation} ${ProssecedElements} ${TableSize}
$ ./go-kmv < ../data/bible.txt
33938 824036 465

# If we (really) count them
$ tr ' ' '\n' < ../data/bible.txt | sort | uniq -c | wc -l
34040
```

If what you want is to use it as a dependency for your project

```go
package main

import gokmv "github.com/positiveblue/go-kmv"

func main() {
    // Get dataStream
    dataStream := myDataStream()

    // Create the estimator
    initialSize := 64 
    estimator := gokmv.NewKMV(initialSize)
    for element := range dataStream {
        // element has to be a UInt64
        estimator.InsertUint64(element)
    }

    estimator.Size() // returns the table size
    estimator.ElementsAdded() // returns the total elements that we processed
    estimator.EstimateCardinality() // returns the cardinality estimation
}
```

Because of the lack of generics in Go go-kmv only provides `Insert` functions for `Uint64` and `strings`. If you want to use your own hash functions or add new types you can just create your own function:

```go
// Insert my type to the table
// Using my hash function
func (kmv *KMV) InsertMyType(s string) {
    // Remember to use the internal seed to have reproducible results
	hash := myHashFunction.Sum64([]byte(s), kmv.Seed())
    // The has has to return a Uint64
	kmv.InsertUint64(hash)
}
```

The formula used for estimating the cardinality is exactly the same described in the paper [ Counting distinct elements in a data stream](http://www.google.com/url?sa=t&rct=j&q=&esrc=s&source=web&cd=1&ved=0CEwQFjAA&url=http%3A%2F%2Fwww.cs.umd.edu%2F~samir%2F498%2Fdistinct.ps&ei=h-3IT5GPBfD16AG0q70v&usg=AFQjCNG4nYiSedl6W3r73ZCXNtnaOancnQ&sig2=E8KzKp4qkLiWMQk690Moyw). What makes this implementation interesting is the use of an adaptive table which grows in order to provide better estimations. The implementation of the adaptive-table can be found [here](https://github.com/positiveblue/adaptive-table)

Cardinalty Estimation is considered solved under all meanings. Nowadays computers have enough memory for computing the cardinality of small sets and for extream cases (big data)algorithms like HyperLogLog and KMV already give an accuracy of ~98% using a few bytes of memory. 

In real life what people usually use is an implementation of [HyperLogLog](http://static.googleusercontent.com/external_content/untrusted_dlcp/research.google.com/en/us/pubs/archive/40671.pdf) with a table size from about 128 to 4096. HyperLogLog and all the algorithms of its family can only use tables of size `2^k` where k is a positive integer. **go-kmv** does not have that limitation and automatically provides a good trade-off without knowing in advance the order of distinct elements that we have to estimate.

The current implementation grows with a factor of `klog(n)` where `k` is the inital table size and `n` is the number of disctinct elements in the stream. That means that runing go-kmv with an `initialSize` of 64 and processing and stream of 10^6 elements the final table size will be about ~600 and the accuracy of the estimation will be ~98.00%.

A critical part to achive meaningful results is to use a good hash function (where good = few colisions). Hash Functions like **FNV**, from the go stdlib are not good enough to ensure the theoretical results. Other algorithms like **AES** provide the best results but are slower and it seems a bit overkill for this implementation. [Murmur3](github.com/spaolacci/murmur3) provides the best ratio results/processing time and it has been used in this implementation.

