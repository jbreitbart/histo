package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/buger/goterm"
	"github.com/jbreitbart/gohistogram"
)

var buckets = flag.Int("buck", 40, "maximum number of buckets")
var updateEvery = flag.Int("update", -1, "print histogram every values")
var filename = flag.String("file", "/dev/stdin", "filename to the files with the cycle")
var weighted = flag.Bool("weighted", false, "Use a weighted histogram that slowly fades out old values")

func main() {
	flag.Parse()

	file, err := os.Open(*filename)
	if err != nil {
		log.Fatalf("Error opening file: %v\n", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var histo gohistogram.Histogram

	if *weighted {
		histo = gohistogram.NewWeightedHistogram(*buckets, 0.064516129)
	} else {
		histo = gohistogram.NewHistogram(*buckets)
	}

	// loops until end or error
	counter := 1
	for scanner.Scan() {
		s := scanner.Text()
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			log.Fatalf("Could not convert file entry (%v) to number: %v\n", s, err)
		}
		histo.Add(n)
		if *updateEvery != -1 && counter%*updateEvery == 0 {
			goterm.Clear()
			goterm.Println(histo.String())
			goterm.Flush()
			goterm.MoveCursor(1, 1)
			time.Sleep(time.Second)
		}
		counter++
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error while scanning the file %v: %v\n", filename, err)
	}

	fmt.Println(histo.String())
}
