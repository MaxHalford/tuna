# tuna

`tuna` is a simple library for computing machine learning features in an online manner. Sometimes datasets are rather large and it isn't convenient to handle them in memory. One approach is to compute running statistics that provide a good approximation of their batch counterparts. The goal of `tuna` is to cover common use cases (*e.g.* a group by followed by a mean) while also making it easy to build custom features.

:warning: I'm working on this for an ongoing Kaggle competition, things are still in flux and the documentation isn't finished

## Quickstart

## API

:point_up: Please check out the [godoc page](https://godoc.org/github.com/MaxHalford/tuna) in addition to the following documentation.

### Extractors

#### `Mean`

The `Mean` struct computes an approximate average. For every `value` the update formula is `mean = mean + (value - mean) / n`. For convenience you can instantiate a `Mean` with the `NewMean` method.

#### Writing a custom feature extractor

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

### Using `GroupBy`

Computing running statistics

### Using `Union`

### Streams

#### `RowStream`

#### `CSVStream`

#### `StreamZip`

The `StreamZip` struct can be used to stream over multiple files without having to concatenate them. Indeed in practice large datasets are more often than not split into chunks for practical reasons. The issue is that if you're using a `GroupBy` and that the group keys are scattered accross multiple files then processing each file individually won't produce the correct result.

To use a `StreamZip` you simply have to instantiate it with a `Stream` slice. Calling `Next` will iterate over each `Row` of each `Stream` and then stop once each `Stream` is depleted.

```go
cs, _ := tuna.NewCSVStream("path/to/file.csv")

sz := tuna.StreamZip{[]tuna.Stream{
    tuna.RowStream{[]tuna.Row{
        tuna.Row{"x0": "42.42", "x1": "24.24"},
        tuna.Row{"x0": "13.37", "x1": "31.73"},
    }},
    cs
}}
```

## Roadmap

- Unit tests
- [Running median](https://rhettinger.wordpress.com/tag/running-median/) (and quantiles!)
- DSL
- CLI tool based on the DSL
- Cute logo

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
