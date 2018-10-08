<div align="center">
  <!-- Logo -->
  <img src="https://docs.google.com/drawings/d/e/2PACX-1vRAWmJOWkS7IByWDZCJQqZmyp2-LO7VPWGgxb9OfuLLFLiquasU3NrS132JyvzkoOx9HcM5DPY2V1-B/pub?w=412&amp;h=213" alt="logo"/>
</div>

<div align="center">
  <!-- godoc -->
  <a href="https://godoc.org/github.com/MaxHalford/tuna">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square" alt="godoc" />
  </a>
  <!-- License -->
  <a href="https://opensource.org/licenses/MIT">
    <img src="http://img.shields.io/:license-mit-ff69b4.svg?style=flat-square" alt="license"/>
  </a>
  <!-- Build status -->
  <a href="https://travis-ci.org/MaxHalford/tuna">
    <img src="https://img.shields.io/travis/MaxHalford/tuna/master.svg?style=flat-square" alt="build_status" />
  </a>
  <!-- Test coverage -->
  <a href="https://coveralls.io/github/MaxHalford/tuna?branch=master">
    <img src="https://coveralls.io/repos/github/MaxHalford/tuna/badge.svg?branch=master&style=flat-square" alt="test_coverage" />
  </a>
</div>

<br/>
<br/>

:warning: I'm working on this for an ongoing Kaggle competition, things are still in flux and the documentation isn't finished

`tuna` is a simple library for computing machine learning features in an online manner. In other words, `tuna` is a streaming ETL. Sometimes datasets are rather large and it isn't convenient to handle them in memory. One approach is to compute running statistics that provide a good approximation of their batch counterparts. The goal of `tuna` is to cover common use cases (*e.g.* a group by followed by a mean) while keeping it simple to build custom features.

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
    // For the purpose of this example we inline the data
    in := `name,£,bangers
"Del Boy",-42,1
Rodney,1001,1
Rodney,1002,2
"Del Boy",42,0
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
    sink, _ := NewCSVSink(os.Stdout)

    // Run
    Run(stream, extractor, sink, 0)
}
```

Running this script will produce the following output in your terminal:

```csv
bangers_sum,name,£_mean
1,Del Boy,0
3,Grandad,0
3,Rodney,1001.5
```

## API

:point_up: Please check out the [godoc page](https://godoc.org/github.com/MaxHalford/tuna) in addition to the following documentation.

### Streams

#### `RowStream`

#### `CSVStream`

#### `StreamZip`

The `StreamZip` struct can be used to stream over multiple files without having to concatenate them. Indeed in practice large datasets are more often than not split into chunks for practical reasons. The issue is that if you're using a `GroupBy` and that the group keys are scattered accross multiple files then processing each file individually won't produce the correct result.

To use a `StreamZip` you simply have to instantiate it with one or more `Stream`s. Calling `Next` will iterate over each `Row` of each `Stream` and then stop once each `Stream` is depleted. Naturally you can combine different types of `Stream`s.

```go
cs, _ := tuna.StreamCSV("path/to/file.csv")

sz := tuna.StreamZip{
    tuna.NewStream(
        tuna.Row{"x0": "42.42", "x1": "24.24"},
        tuna.Row{"x0": "13.37", "x1": "31.73"},
    ),
    cs
}
```

#### Writing a custom `Stream`

### Extractors

#### `Mean`

The `Mean` struct computes an approximate average. For every `value` the update formula is `mean = mean + (value - mean) / n`. For convenience you can instantiate a `Mean` with the `NewMean` method.

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

The `Collect` methods returns a channel that streams `Row`s. Each such `Row` will then be stored in a CSV file (depending on your application). Most `Extractor`s only return a single result. For example `Mean` returns a `Row` with one key named `"mean"` and one value representing the current mean. On the other `GroupBy` returns multiple `Row`s (one per group).

The `Size` method is simply here to monitor the number of computed values. Most `Extractor`s simply return 1 whereas `GroupBy` returns the sum of the sizes of each group.

Naturally the easiest way to proceed is to copy/paste one of the existing `Extractor`s and then edit it.

### `GroupBy`

### `Union`

### Sinks

#### `CSVSink`

#### Writing a custom `Sink`

### The `Run` method

## Roadmap

- Unit tests
- [Running median](https://rhettinger.wordpress.com/tag/running-median/) (and quantiles!)
- DSL
- CLI tool based on the DSL
- Handle dependencies between extractors (for example `Variance` could reuse `Mean`)
- Identify bottlenecks

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
