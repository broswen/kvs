package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

var hostFlag = flag.String("host", "", "host to call")
var numFlag = flag.Int("n", 1, "number of times to call")

func main() {
	rand.Seed(time.Now().Unix())
	flag.Parse()
	if *hostFlag == "" {
		log.Fatalf("Must specify host flag.")
	}

	batchBenchmark("Set", testSet(*hostFlag), *numFlag)
	batchBenchmark("Get", testGet(*hostFlag), *numFlag)
}

// get a random item from the host, this also takes time to generate random key
func testGet(host string) func() {
	return func() {
		key := strconv.Itoa(rand.Intn(1000))
		_, err := http.Get(fmt.Sprintf("%s/%s", host, key))
		if err != nil {
			log.Printf("get: %v\n", err)
		}
	}
}

// set a random item from the host, this also takes time to generate random key/value
func testSet(host string) func() {
	return func() {
		key := strconv.Itoa(rand.Intn(1000))
		value := strconv.Itoa(rand.Intn(1000))
		_, err := http.Post(fmt.Sprintf("%s/%s", host, key), "text/plain", bytes.NewBufferString(value))
		if err != nil {
			log.Printf("get: %v\n", err)
		}
	}
}

// run a function and return the duration
func benchmark(fn func()) time.Duration {
	start := time.Now()
	fn()
	return time.Since(start)
}

// benchmark a function n times and print the average
func batchBenchmark(name string, fn func(), n int) {
	times := make([]int64, 0)

	for i := 0; i < *numFlag; i++ {
		time := benchmark(fn)
		times = append(times, int64(time))
	}

	average := average(times)

	log.Printf("average of %v for %s\n", time.Duration(average), name)

}

// find the average for a slice of int64
func average(nums []int64) int64 {
	var total int64 = 0
	for i := 0; i < len(nums); i++ {
		total += nums[i]
	}
	return total / int64(len(nums))
}
