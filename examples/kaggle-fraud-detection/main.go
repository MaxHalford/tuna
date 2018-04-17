package main

import (
	"fmt"
	"hash/fnv"
	"log"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/MaxHalford/line"
)

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {

	var (
		nFeatures = uint32(math.Pow(2, 26))
		trainPath = "train.csv"
		testPath  = "test.csv"
		subPath   = "submission_ftrl.csv"
		nEpochs   = 1
	)

	// Choose what combinations of features to one-hot encode
	var (
		order1 = []string{"ip", "app", "device", "os", "channel", "click_hour"}
		order2 = [][]string{
			[]string{"ip", "app"},
			[]string{"ip", "channel"},
			[]string{"ip", "device"},
			[]string{"ip", "os"},
			[]string{"app", "channel"},
		}
		order3 = [][]string{
			[]string{"ip", "app", "os"},
		}
		order4 = [][]string{
			[]string{"ip", "device", "app", "os"},
		}
	)

	// Store the last click time per IP
	var lastClickTimes = make(map[string]time.Time)

	// Define a RowParser
	var rp = func(row line.Row) (id string, x line.Vector, y float64) {
		x = make(line.Vector)

		// Time since last click per IP
		clickTime, _ := time.Parse(time.RFC3339, row["click_time"])
		lastClickTime, ok := lastClickTimes[row["ip"]]
		if ok {
			x[nFeatures-1] = float64(clickTime.Sub(lastClickTime))
		} else {
			x[nFeatures-1] = -1
		}
		lastClickTimes[row["ip"]] = clickTime

		// Click hour
		row["click_hour"] = row["click_time"][11:13]

		// Order 1 features
		for _, name := range order1 {
			x[hash(name+"_"+row[name])%(nFeatures-1)]++
		}

		// Order 2 features
		for _, pair := range order2 {
			x[hash(pair[0]+"_"+row[pair[0]]+pair[1]+"_"+row[pair[1]])%(nFeatures-1)]++
		}

		// Order 3 features
		for _, pair := range order3 {
			x[hash(pair[0]+"_"+row[pair[0]]+pair[1]+"_"+row[pair[1]]+pair[2]+"_"+row[pair[2]])%(nFeatures-1)]++
		}

		// Order 4 features
		for _, pair := range order4 {
			x[hash(pair[0]+"_"+row[pair[0]]+pair[1]+"_"+row[pair[1]]+pair[2]+"_"+row[pair[2]]+pair[3]+"_"+row[pair[3]])%(nFeatures-1)]++
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
	var model = line.NewFTRLProximalClassifier(0.2, 0.8, 0.01, 0)

	// Train
	for i := 1; i <= nEpochs; i++ {
		fmt.Printf("Epoch %d\n", i)
		ri, err := line.NewCSVRowReader(trainPath)
		if err != nil {
			log.Fatal(err)
		}
		line.Train(model, ri, rp, line.Accuracy{}, os.Stdout, 1000000)
	}

	// Submission
	sub, err := os.Create(subPath)
	if err != nil {
		log.Fatal(err)
	}
	sub.WriteString("click_id,is_attributed\n")
	ri, err := line.NewCSVRowReader(testPath)
	if err != nil {
		log.Fatal(err)
	}
	line.Predict(model, ri, rp, sub, os.Stdout, 1000000)
}
