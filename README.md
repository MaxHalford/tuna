# line

`line` is a simple library for doing [online learning](https://www.wikiwand.com/en/Online_machine_learning). In short the goal of online learning is to train a machine learning model one observation at a time. With `line` the user total freedom as to how to preprocess the data. `line` takes over once the user has provided a `RowReader` to stream rows and a `RowParser` to parse them.

:warning: For the while only binary classification is supported.

## Usage

Check out [GoDoc](https://godoc.org/github.com/MaxHalford/gago) for a full overview.

- A `Row` is a `map[string]string` that is streamed from somewhere (for example a CSV file)
- A `RowReader` implements a `Read()` method which returns a `Row`
- An `Instance` is a parsed version of a `Row`
- A `RowParser` is a method that parses a `Row` into an `Instance`
- A `Model` implements `FitPartial` and `PredictPartial`
- A `Metric` computes an online metric

The idea is that the user only has to implement a `RowReader` and a `RowParser`. The `RowReader` can be anything as long as it implements a `Read` method that returns a `Row`. For example `line` provides a `CSVRowReader` to stream from a CSV file row by row. The `RowParser` has the following signature:

```go
func(row Row) (ID string, x Vector, y float64)
```

It has to take as input a `Row` and output a row identifier, a `Vector`, and a target. A `Vector` has the following signature:

```go
type Vector map[uint32]float64
```

A `Vector` thus maps unsigned integers to floating point values. This allows a sparse representation which is commonly needed for large-scale online learning. The user has the responsibility (and freedom) to decide what features to include. Typically features should be hashed to obtain a `uint32` value.

Once a `RowReader` and a `RowParser` have been defined the user can pick a `Model` and call the `Train` method which as the following signature:

```go
func Train(model Model, ri RowReader, rp RowParser, metric Metric, monitor *os.File, monitorEvery uint64)
```

The metric is used to monitor the performance of the `Model` every `monitorEvery` instances. The metric is applied to each instance *before* it is fed to the model. The output is piped via the `monitor` argument (`os.Stdout` can be used in most cases).

Once training has terminated, the `Predict` method can be used to make predictions on another stream of instances:

```go
func Predict(model Model, ri RowReader, rp RowParser, output *os.File, monitor *os.File, monitorEvery uint64)
```

The `ID` and the prediction of each instance will be written to the `output` file (typically you can use `os.Create`). Progress is sent to the `monitor` file.

## Models

### Follow The Regularised Leader Proximal (FTRL-Proximal)

FTRL-Proximal is a logistic regression with an adaptive learning rate and L1-L2 regularisation.

```go
model := line.NewFTRLProximalClassifier(0.2, 0.8, 0.01, 0) // alpha, beta, l1 regularisation, l2 regularisation
```

- [Paper](http://www.eecs.tufts.edu/~dsculley/papers/ad-click-prediction.pdf)
- [Example](examples/kaggle-fraud-detection)

## Metrics

- `line.Accuracy`
- `line.LogLoss`

## Row readers

- `line.CSVRowReader`

## Dependencies

None.

## To do

- Perceptron
- Passive-Aggressive
- Winnow
- Mondrian trees
    - http://papers.nips.cc/paper/3622-the-mondrian-process.pdf
    - https://arxiv.org/pdf/1406.2673.pdf
    - https://scikit-garden.github.io/examples/MondrianTreeRegressor/
- FFM

## License

The MIT License (MIT). Please see the [LICENSE file](LICENSE.md) for more information.
