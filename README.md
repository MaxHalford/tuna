<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRAWmJOWkS7IByWDZCJQqZmyp2-LO7VPWGgxb9OfuLLFLiquasU3NrS132JyvzkoOx9HcM5DPY2V1-B/pub?w=412&amp;h=213" alt="logo"/>
</div>

<div align="center">
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/tuna">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/tuna">
    <img src="https://img.shields.io/travis/MaxHalford/tuna/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Test coverage -->
  <a href="https://coveralls.io/github/MaxHalford/tuna?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/tuna/badge.svg?branch=master&style=flat-square" alt="test_coverage" />
  </a>
  <!-- Go report card -->
  <a href="https://goreportcard.com/report/github.com/MaxHalford/tuna">
    <img src="https://goreportcard.com/badge/github.com/MaxHalford/tuna?style=flat-square" alt="go_report_card" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
</div>

<br/>
<br/>


`tuna` is a framework for computing streaming aggregates. In other words `tuna` is a streaming ETL. Sometimes datasets don't fit in memory and so you have to process them in chunks. One approach is to compute running statistics that provide a good approximation of their batch counterparts. The goal of `tuna` is to cover common use cases (e.g. a group by followed by a mean) while keeping it simple to build custom features.


## Concepts

Like many [such libraries](https://github.com/topics/etl), `tuna` involves a few simple concepts:

- A `Row` is a set of (key, value) pairs (represented in Go with a `map[string]string`)
- A `Stream` is a source of data that returns `Row`s one by one
- A `Metric` is an object that computes one or more running statistics; it is fed `float64` values and returns a `map[string]float64` of features
- An `Agg` takes `Row`s in, extracts `float64`s, and passes them to one or more `Metric`s
- A `Sink` persists the output of an `Agg`


## Quickstart

```go
package main

import (
    "os"
    "strings"

    "github.com/MaxHalford/tuna"
)

func main() {
    // For the sake of example we inline the data, but usually it should be
    // located in a file, database, or some other source
    in := `name,£,bangers
Del Boy,-42,1
Rodney,1001,1
Rodney,1002,2
Del Boy,42,0
Grandad,0,3`

    // Define a Stream
    stream, _ := NewCSVStream(strings.NewReader(in))

    // Define an Agg
    agg := NewGroupBy(
        "name",
        func() Agg {
            return Aggs{
                NewExtractor("£", NewMean(), NewSum()),
                NewExtractor("bangers", NewSum()),
            }
        },
    )

    // Define a Sink
    sink := NewCSVSink(os.Stdout)

    // Run
    Run(stream, agg, sink, 0) // 0 means we don't display live progress
}
```

Running this script will produce the following output in your terminal:

```csv
bangers_sum,name,£_mean,£_sum
1,Del Boy,0,0
3,Grandad,0,0
3,Rodney,1001.5,2003
```


## Usage

:point_up: In addition to the following documentation, please check out the [godoc page](https://godoc.org/github.com/MaxHalford/tuna) for detailed information.

### Streams

#### Streaming from a CSV file

A common use case you may have is processing rows located in a CSV file. You can use the `NewCSVStream` method to stream CSV data from an `io.Reader` instance.

```go
var r io.Reader // Depends on your application
stream := tuna.NewCSVStream(r)
```

For convenience you can use the `NewCSVStreamFromPath` method to stream CSV data from a file path, it is simply a wrapper on top of `NewCSVStream`.

```go
stream := tuna.NewCSVStreamFromPath("path/to/file")
```

#### Streaming `Rows` directly

For some reason you may want to stream a given set of `Row`s. Although this defeats the basic paradigm of `tuna` which is to process data that can't be loaded in memory, it is practical for testing purposes.

```go
stream := tuna.NewStream(
    tuna.Row{"x0": "42.42", "x1": "24.24"},
    tuna.Row{"x0": "13.37", "x1": "31.73"},
)
```

#### Streaming from multiple sources

The `ZipStreams` method can be used to stream over multiple sources without having to merge them manually. Indeed large datasets are more often than not split into chunks for practical reasons. The issue is that if you're using a `GroupBy` and that the group keys are scattered across multiple sources, then processing each file individually won't produce the correct result.

To use `ZipStreams` you simply have to provide it with one or more `Stream`s. It will then return a new `Stream` which will iterate over each `Row` of each provided `Stream` until they are all depleted. Naturally you can combine different types of `Stream`s.

```go
s1, _ := tuna.NewCSVStreamFromPath("path/to/file.csv")
s2 := tuna.NewStream(
    tuna.Row{"x0": "42.42", "x1": "24.24"},
    tuna.Row{"x0": "13.37", "x1": "31.73"},
)
stream := tuna.ZipStreams(s1, s2)
```

#### Using a custom `Stream`

A `Stream` is simply a channel that returns `ErrRow`s, i.e.

```go
type Stream chan ErrRow
```

An `ErrRow` has the following signature.

```go
type ErrRow struct {
    Row
    Err error
}
```

A `Row` is nothing more than a `map[string]string`. The `Err` fields indicates if something went wrong during the retrieval of the corresponding `Row`.

### Metrics

#### Overview

`Metric`s are the objects that do the actual computation. You're supposed to use them by providing them to an `Extractor`. Every time the `Extractor` is fed a `Row`, it will simply parse the values from the `"banana"` field and feed them to each `Metric`.

The current list of available metrics are:

- `Mean`
- `Variance`
- `Sum`
- `Min`
- `Max`
- `Skew`
- `Kurtosis`
- `Diff`

Although you can instantiate each struct yourself, it is recommended that you instantiate each struct with it's respective `New` method. For example use the `NewMax` method if you want to use the `Max` struct.

:point_up: A set of `Metric`s, which is represented in `tuna` by the `Metrics`, is also a `Metric`.

#### Writing a custom `Metric`

Every metric has to implement the following interface:

```go
type Metric interface {
    Update(x float64) error
    Collect() map[string]float64
}
```

- The `Update` method updates the running statistic that is being computed. For example the update formula for the `Mean` metric is `mean = mean + (x - mean) / n`.
- The `Collect` method returns a set of one or more features. For example the `Mean` metric returns `{"mean": some_float64}`.


### Aggs

#### Overview

`Agg`s are the bridge between `Row`s and `Metric`s. The simplest type of `Agg` is the `Extractor`, which extracts a `float64` value from a `Row` and feeds to a `Metric`. Another example is the `GroupBy` struct, which maintains a set of set of `Agg`s and feeds them values given a `Row` key. `Agg`s can be composed to build powerful and expressive pipelines.

#### Extractor

As already said the `Extractor` is the simplest kind of `Agg`. It has the following signature:

```go
type Extractor struct {
    Extract func(row Row) (float64, error)
    Metric  Metric
    Prefix  string
}
```

Simply put an `Extractor`s parses a `Row` and extracts a `float64` using it's `Extract` method. It then feeds the `float64` to the `Metric`. After retrieving the results by calling the `Metric`'s `Collect` method the `Extractor` will prepend the `Prefix` to each key so as to add the field name to the results.

The `Extract` field gives you the flexibility of parsing each `Row` as you wish. However often you might simply want to cast each value as `float64`. In this case you can use the `NewExtractor` method for convenience, as so:

```go
extractor := tuna.NewExtractor("banana", tuna.NewMean(), tuna.NewMedian())
```

#### `GroupBy`

Computing running statistics is nice but in practice you probably want to compute conditional statistics. In other words you want to "group" the incoming values by a given attribute and compute one or more statistics inside each group. This is what the `GroupBy` struct is intended for.

You can use the `NewGroupBy` method to instantiate a `GroupBy`, it takes as arguments a `string` which tells it by what field to group the data and a `func() Agg` callable which returns an `Agg`. Every time a new key appears the callable will be used to instantiate a new `Agg` for the new group.

```go
gb := tuna.NewGroupBy("name", func() tuna.Agg { return tuna.NewExtractor("£", tuna.NewSum()) })
```

You can nest `GroupBy`s if you want to group the data by more than one variable. For example the following `Agg` will count the number of taken bikes along with the number of returned bikes by city as well as by day.

```go
gb := tuna.NewGroupBy(
    "city",
    func() tuna.Agg {
        return tuna.NewGroupBy(
            "day",
            func() tuna.Agg {
                return tuna.Aggs{
                    tuna.NewExtractor("bike_taken", tuna.NewSum()),
                    tuna.NewExtractor("bike_returned", tuna.NewSum()),
                )
            }
        )
    }
)
```

#### `SequentialGroupBy`

Using a `GroupBy` can incur a large memory usage if you are computing many statistics on a very large dataset. Indeed the spatial complexity is `O(n * k)`, where `n` is the number of group keys and `k` is the number of `Agg`s. This can potentially become quite large, especially if you're using nested `GroupBy`s. While this is completely manageable if you have enough available RAM, it can still hinder the overall computation time.

The trick is that **if your data is ordered by the group key then you only have to store the running statistics for one group at a time**. This leads to an `O(k)` spatial complexity which is much more efficient. While having ordered data isn't always the case, you should make the most of it if it is. To do so you can use the `SequentialGroupBy` struct which can be initialized with the `NewSequentialGroupBy` method. It takes as argument a `Sink` in addition to the arguments used for the `NewGroupBy` method. Every time a new group key is encountered the current statistics are flushed to the `Sink` and a new `Agg` is initialized to handle the new group.

```go
stream, _ := tuna.NewCSVStreamFromPath("path/to/csv/ordered/by/name")

sink, _ := tuna.NewCSVSinkFromPath("path/to/sink")

sgb := tuna.NewSequentialGroupBy(
    "name",
    func() tuna.Agg { return tuna.NewExtractor("bangers", NewMean()) }
    sink
)

tuna.Run(stream, sgb, nil, 1e6)
```

:point_up: Make sure your data is ordered by the group key before using `SequentialGroupBy`. There are various ways to sort a file by a given field, one of them being the [Unix `sort` command](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/sort.html).


### Sinks

#### `CSVSink`

You can use a `CSVSink` struct to write the results of an `Agg` to a CSV file. It will write one line for each `Row` returned by the `Agg`'s `Collect` method. Use the `NewCSVSink` method to instantiate a `CSVSink` that writes to a given `io.Writer`.

```go
var w io.Reader // Depends on your application
sink := tuna.NewCSVSink(r)
```

For convenience you can use the `NewCSVStreamFromPath` method to stream CSV data from a file path, it is simply a wrapper on top of `NewCSVStream`.

```go
sink := tuna.NewCSVSinkFromPath("path/to/file")
```

#### Writing a custom `Sink`

The `Sink` interface has the following signature:

```go
type Sink interface {
    Write(rows <-chan Row) error
}
```

A `Sink` simply has to be able to write a channel of `Row`s "somewhere".


### The `Run` method

Using the `Run` method is quite straightforward.

```go
checkpoint := 1e5
err := Run(stream, agg, sink, checkpoint)
```

You simply have to provide it with a `Stream`, an `Agg`, and a `Sink`. It will feed the `Agg` with the `Row`s produced by the `Stream` one by one. Once the `Stream` is depleted the results of the `Agg` will be written to the `Sink`. An `error` will be returned if anything goes wrong along the way. The `Run` method will also display live progress in the console every time the number of parsed rows is a multiple of `checkpoint`, e.g.

```sh
00:00:02 -- 300,000 rows -- 179,317 rows/second
```

:point_up: In the future there might be a `Runner` interface to allow more flexibility. In the meantime you can copy/paste the content of the `Run` method and modify it as needed if you want to do something fancy (like monitoring progress inside a web page or whatnot)

## Roadmap

- [Running median](https://rhettinger.wordpress.com/tag/running-median/) (and quantiles!)
- DSL
- CLI tool based on the DSL
- Maybe handle dependencies between aggs (for example `Variance` could reuse `Mean`)
- Benchmark and identify bottlenecks

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
