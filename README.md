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

:warning: I'm working on this for an ongoing Kaggle competition, things are still in flux and the documentation isn't finished

`tuna` is a simple library for computing machine learning features in an online manner. In other words, `tuna` is a streaming ETL. Sometimes datasets are rather large and it isn't convenient to handle them in memory. One approach is to compute running statistics that provide a good approximation of their batch counterparts. The goal of `tuna` is to cover common use cases (e.g. a group by followed by a mean) while keeping it simple to build custom features.

Like many [such libraries](https://github.com/topics/etl), `tuna` involves a few simple concepts:

- A `Row` is a set of key:value pairs (represented in Go with a `map[string]string`)
- A `Stream` is a source of data that returns `Row`s one by one
- An `Extractor` is fed `Rows` one by one and maintains one or more running statistics
- A `Sink` persists the output of an `Extractor`

You can then use the `Run` method to stitch these together.

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

    // Define an Extractor
    extractor := NewGroupBy(
        "name",
        func() Extractor {
            return NewUnion(
                NewMean("£"),
                NewSum("bangers"),
            )
        },
    )

    // Define a Sink
    sink := NewCSVSink(os.Stdout)

    // Run
    Run(stream, extractor, sink, 0) // 0 means we don't display live progress
}
```

Running this script will produce the following output in your terminal:

```csv
bangers_sum,name,£_mean
1,Del Boy,0
3,Grandad,0
3,Rodney,1001.5
```

## Usage

:point_up: Please check out the [godoc page](https://godoc.org/github.com/MaxHalford/tuna) in addition to the following documentation.

### Streams

#### Streaming from a CSV file

The most common use case you may probably have is to process rows located in a CSV file. You can use `NewCSVStream` to stream CSV data from an `io.Reader` instance.

```go
var r io.Reader // Depends on your application
stream := tuna.NewCSVStream(r)
```

For convenicence you can use the `NewCSVStreamFromPath` method to stream CSV data from a file path, it is simply a wrapper on top of `NewCSVStream`.

```go
stream := tuna.NewCSVStreamFromPath("path/to/file")
```

#### Streaming `Rows` directly

For some reason you might want to stream from a given set of `Row`s. However this defeats the basic paradigm of `tuna` which is that the data can't be loaded in memory in it's entirety. Regardless streaming `Row`s directly is practical for testing purposes.

```go
stream := tuna.NewStream(
    tuna.Row{"x0": "42.42", "x1": "24.24"},
    tuna.Row{"x0": "13.37", "x1": "31.73"},
)
```

#### Streaming from multiple sources

The `ZipStreams` method can be used to stream over multiple sources without having to concatenate them. Indeed in practice large datasets are more often than not split into chunks for practical reasons. The issue is that if you're using a `GroupBy` and that the group keys are scattered accross multiple sources, then processing each file individually won't produce the correct result.

To use `ZipStreams` you simply have to provide it with one or more `Stream`s. It will then return a `Stream` which will iterate over each `Row` of each provided `Stream` until they are all depleted. Naturally you can combine different types of `Stream`s.

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

### Extractors

#### Basic aggregators

All of the following `Extractor`s work in the same way:

- `Mean`
- `Variance`
- `Sum`
- `Min`
- `Max`
- `PTP` which stands for "peak to peak" and computes `max(x) - min(x)`
- `Skew`
- `Kurtosis`

For convenience you can instantiate each struct with it's respective `New` method. For example use the `NewMean` method if you want to use the `Mean` struct. Each of these methods takes as argument a `string` which indicates the field for which the aggregate should be computed.

You can have finer control by modifying each struct after calling it's `New` method. Each of the above structs has a `Parse` method and a `Prefix` string. For example the signature of the `Max` struct is:

```go
type Max struct {
    Parse  func(Row) (float64, error)
    Prefix string
}
```

The `Parse` method determines how to parse an incoming `Row`. This allows doing fancy things, for example parsing a field as a `float64` and then applying a logarithmic transform. The `Prefix` string determines what prefix should be used when the results are collected. For example setting a `candy_` prefix with the `Max` struct will use `candy_max` as a result name.

:warning: It isn't recommended to instantiate an `Extractor` yourself. Calling the `New` methods sets initial values which are required for obtaining correct results. If you want to modify the `Parse` or the `Prefix` fields then set them after calling the `New` method.

#### `Diff`

A common use case that occurs for ordered data is to compute statistics of the differences between consecutive values. This requires memorising the previous value and feeding the difference with the current value to an `Extractor`. You can do this by using the `Diff` struct which has the following signature:

```go
type Diff struct {
    Parse     func(Row) (float64, error)
    Extractor Extractor
    FieldName string
}
```

The `Parse` method tells the `Diff` how to parse an incoming `Row`. The current value will be stored and update each time a new `Row` comes in. When this happens a new field will be set on the `Row`, which will be then be processed through the `Extractor`.  The `FieldName` string determines the name of this new field. It's important **that the `Extractor` parsed the field named `FieldName`**. For convenience you can use the `NewDiff` method if you don't have to do any fancy processing.

:point_up: Make sure your data is ordered in the right way before using `Diff`. There are various ways to sort a file by a given field, one of them being the [Unix `sort` command](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/sort.html).

#### `Union`

Most use cases usually involve computing multiple statistics. One way would be to define a single `Extractor` which computes several statistics simultaneously. While this is computationally efficient and allows reusing computed values, it leads to writing application specific code that is somwhat difficult to maintain.

Another way is to define a slice of `Extractors`s and loop over each one of them every time a new `Row` comes in. This is exactly what the `Union` struct does. You can use the [variadic](https://gobyexample.com/variadic-functions) `NewUnion` method to instantiate a `Union`.

```go
union := NewUnion(NewMean("flux"), NewSum("flux"))
```

You can think of this as [composition](https://www.wikiwand.com/en/Function_composition_(computer_science)).

#### `GroupBy`

Computing running statistics is nice but in practice you probably want to compute conditional statistics. In other words you want to "group by" the incoming values by a given variable and compute one or more statistics inside each group. This is what the `GroupBy` struct is intended for.

You can use the `NewGroupBy` method to instantiate a `GroupBy`, it takes as arguments a `string` which tells it by what field to group the data and a `func() Extractor` callable which returns an `Extractor`. Every time a new key appears the callable will be used to instantiate a new `Extractor` for the new group.

```go
gb := NewGroupBy("name", func() Estimator { return NewSum("£") })
```

A nice thing is that you can use `Union`s together with `GroupBy`s. This will compute multiple statistics per group.

```go
gb := NewGroupBy(
    "name",
    func() Extractor {
        return NewUnion(
            NewMean("£"),
            NewSum("bangers"),
        )
    },
)
```

You can nest `GroupBy`s if you want to group the data by more than one variable. For example the following `Extractor` will count the number of taken bikes along with the number of returned bikes by city as well as by day.

```go
gb := NewGroupBy(
    "city",
    func() Estimator {
        return NewGroupBy(
            "day",
            func() Estimator {
                return NewUnion(
                    NewSum("bike_taken"),
                    NewSum("bike_returned"),
                )
            }
        )
    }
)
```

#### `SequentialGroupBy`

Using a `GroupBy` can incur a large memory usage if you are computing many statistics on a very large dataset. Indeed the spatial complexity is `O(n * k)`, where `n` is the number of group keys and `k` is the number of `Extractor`s. This can potentially become quite large, especially if you're using nested `GroupBy`s. While this is completely fine if you have enough RAM available, it can hinder the overall computation time.

The trick is that **if your data is ordered by the group key then you only have to store the running statistics for one group at a time**. This leads to a `O(k)` spatial complexity which is much more efficient. While having ordered data isn't always the case, you should make the most of it is the case. To do so you can use the `SequentialGroupBy` struct which can be initialised with the `NewSequentialGroupBy` method. It takes as argument a `Sink` in addition to the arguments used for the `NewGroupBy` method. Every time a new group key is encountered the current statistics are written to the `Sink` and a new `Extractor` is initialised to handle the new group.

```go
stream, _ := NewCSVStreamFromPath("path/to/csv/ordered/by/name")

sink, _ := tuna.NewCSVSinkFromPath("path/to/sink")

sgb := tuna.NewSequentialGroupBy(
    "name",
    func() Extractor { return NewMean("bangers") }
    sink
)

tuna.Run(stream, sgb, nil, 1e6)
```

:point_up: If you're using a `SequentialGroupBy` then you don't have to provide a `Sink` to the `Run` method. This is because the results will be written every time a new group key is encountered.

:point_up: Make sure your data is ordered by the group key before using `SequentialGroupBy`. There are various ways to sort a file by a given field, one of them being the [Unix `sort` command](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/sort.html).

#### Writing a custom `Extractor`

A feature extractor has to implement the following interface.

```go
type Extractor interface {
    Update(Row) error
    Collect() <-chan Row
    Size() uint
}
```

The `Update` method updates the running statistic that is being computed.

The `Collect` method returns a channel that streams `Row`s. Each such `Row` will be persisted with a `Sink`. Most `Extractor`s only return a single `Row`. For example `Mean` returns a `Row` with one key named `"mean_<field>"` and one value representing the estimated mean. `GroupBy`, however, returns multiple `Row`s (one per group).

The `Size` method is simply here to monitor the number of computed values. Most `Extractor`s simply return `1` whereas `GroupBy` returns the sum of the sizes of each group.

Naturally the easiest way to proceed is to copy/paste one of the existing `Extractor`s and then edit it.

### Sinks

#### `CSVSink`

You can use a `CSVSink` struct to write the results of an `Extractor` to a CSV file. It will write one row per `Row` returned by the `Extractor`'s `Collect` method. Use the `NewCSVSink` method to instantiate a `CSVSink` that writes to a given `io.Writer`.

```go
var w io.Reader // Depends on your application
sink := tuna.NewCSVSink(r)
```

For convenicence you can use the `NewCSVStreamFromPath` method to stream CSV data from a file path, it is simply a wrapper on top of `NewCSVStream`.

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
err := Run(stream, extractor, sink, checkpoint)
```

You simply have to provide it with a `Stream`, an `Extractor`, and a `Sink`. It will update the `Extractor` with the `Row`s produced by the `Stream` one by one. Once the `Stream` is depleted the results of the `Extractor` will be written to the `Sink`. An `error` will be returned if anything goes wrong along the way. The `Run` method will also display live progress in the console everytime the number of parsed rows is a multiple of `checkpoint`, e.g.:

```sh
00:00:02 -- 300,000 rows -- 179,317 rows/second -- 78 values in memory
```

## Roadmap

- [Running median](https://rhettinger.wordpress.com/tag/running-median/) (and quantiles!)
- DSL
- CLI tool based on the DSL
- Maybe handle dependencies between extractors (for example `Variance` could reuse `Mean`)
- Benchmark and identify bottlenecks

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
