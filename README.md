# tuna

`tuna` is a simple library for computing machine learning features in an online manner. Sometimes datasets can be rather large and it isn't convenient to handle them in memory. The idea here is to cover common use cases (e.g. a group by followed by a mean) while also making it easy to build custom features.

:warning: I'm working on this for an ongoing Kaggle competition, things are still in flux and the documentation isn't finished

## Quickstart

## API

:point_up: Please check out the [godoc page]() in addition to the following documentation.

### Extractors

**Mean**

The `Mean` struct computes an approximate average. For every `value` the update formula is `mean = mean + (value - mean) / n`. For convenience you instantiate a `Mean` with the `NewMean` method.

**Writing a custom feature extractor**

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

Naturally the easiest way to proceed is to copy/paste one of the existing `Extractor`s and edit it.

### Using `GroupBy`

### Using `Union`

### Streams

- `CSVStream` (use `NewCSVStream`)

## Roadmap

- Unit tests
- [Running median](https://rhettinger.wordpress.com/tag/running-median/)
- DSL
- CLI tool based on the DSL

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
