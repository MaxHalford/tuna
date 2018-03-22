package main

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"log"
	"os"
	"strconv"

	"github.com/MaxHalford/koza/metrics"
)

func rocauc(y, yHat []float64) float64 {
	var rocauc, _ = metrics.ROCAUC{}.Apply(y, yHat, nil)
	return rocauc
}

func hash(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func oneHot(s string, n uint64) string {
	h := hash(s)
	return strconv.FormatUint(h%n, 10)
}

func main() {

	var (
		//nFeatures = uint64(math.Pow(2, 24))
		trainPath = "../data/kaggle/train.csv"
		testPath  = "../data/kaggle/test.csv"
		subPath   = "submission_ftrl.csv"
		nEpochs   = 3
	)

	// Define a RowParser
	var rp = func(row Row) (id string, x Vector, y float64) {
		x = make(Vector)

		row["click_hour"] = row["click_time"][11:13]

		// Order 1 features
		for _, name := range []string{"ip", "app", "device", "os", "channel", "click_hour"} {
			x[name+"_"+row[name]]++
		}

		// Order 2 features
		var pairs = [][]string{
			[]string{"ip", "app"},
			[]string{"ip", "device"},
			[]string{"ip", "os"},
			[]string{"ip", "channel"},
			[]string{"ip", "click_hour"},
			[]string{"app", "device"},
			[]string{"app", "os"},
			[]string{"app", "channel"},
			[]string{"app", "click_hour"},
			[]string{"device", "os"},
			[]string{"device", "channel"},
			[]string{"device", "click_hour"},
			[]string{"os", "channel"},
			[]string{"os", "click_hour"},
			[]string{"channel", "click_hour"},
		}
		for _, pair := range pairs {
			x[pair[0]+"_"+row[pair[0]]+"_"+pair[1]+"_"+row[pair[1]]]++
		}

		// Parse target
		yStr, ok := row["is_attributed"]
		if ok {
			y, _ = strconv.ParseFloat(yStr, 64)
		} else {
			y = -1
		}

		// Parse ID
		id, ok = row["click_id"]
		if !ok {
			id = ""
		}

		return id, x, y
	}

	// Instantiate the model
	var learner = NewFTRLProximalClassifier(0.2, 1, 1, 1) // NewFMClassifier(4, 0.1)

	// Train
	for i := 1; i <= nEpochs; i++ {
		fmt.Printf("Epoch %d\n", i)
		ri, err := NewCSVRowReader(trainPath)
		if err != nil {
			log.Fatal(err)
		}
		trainStream := NewInstanceStream(ri, rp)
		Fit(learner, trainStream, rocauc, 1000000)
	}

	// Submission
	f, err := os.Create(subPath)
	if err != nil {
		log.Fatal(err)
	}
	sub := csv.NewWriter(f)
	defer sub.Flush()
	sub.Write([]string{"click_id", "is_attributed"})
	ri, err := NewCSVRowReader(testPath)
	if err != nil {
		log.Fatal(err)
	}
	for instance := range NewInstanceStream(ri, rp) {
		yHat := learner.PredictPartial(instance.x)
		sub.Write([]string{instance.ID, strconv.FormatFloat(yHat, 'f', -1, 64)})
	}
}
