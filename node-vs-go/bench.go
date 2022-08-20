package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

func exit() {
	fmt.Println("USAGE: bench.go <url> <parallel-connections>")
	os.Exit(1)
}

func metrics(endTime, execTime int64) {
	// TODO: Implement circular buffer.
}

func work(url string) {
	_, err := http.Get(url)
	if err != nil {
		// TODO: handle error
	}
}

func worker(stopConsumerChan <-chan bool, jobChan <-chan string, w *sync.WaitGroup) {
	w.Add(1)
	fmt.Println(w)
	count := 0
	for {
		select {
		case url := <-jobChan:
			count++
			work(url)
		case <-stopConsumerChan:
			defer w.Done()
			return
		}
	}
}

func generate(url string, stopChan <-chan os.Signal, jobChan chan string) {
	for {
		select {
		case <-stopChan:
			fmt.Println("Benchmark Stopping..")
			return
		default:
			jobChan <- url
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		exit()
	}

	url := os.Args[1]
	conn, err := strconv.Atoi(os.Args[2])

	if err != nil {
		exit()
	}

	fmt.Println("Benchmark Start")

	stopProducerChan := make(chan os.Signal, 1)
	stopConsumerChan := make(chan bool, conn)
	jobChan := make(chan string, conn)

	signal.Notify(stopProducerChan, syscall.SIGINT, os.Interrupt)

	var workerWG sync.WaitGroup
	for i := 0; i < conn; i++ {
		go worker(stopConsumerChan, jobChan, &workerWG)
	}

	generate(url, stopProducerChan, jobChan)

	for i := 0; i < conn; i++ {
		stopConsumerChan <- true
	}
	workerWG.Wait()

	close(stopProducerChan)
	close(jobChan)
	close(stopConsumerChan)

	fmt.Println("Benchmark Done")
}
