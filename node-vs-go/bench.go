package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

type token struct{}

func exit() {
	fmt.Println("USAGE: bench.go <url> <parallel-connections>")
	os.Exit(1)
}

func metrics(endTime, execTime int64) {
	// TODO: Implement circular buffer.
}

func work(url string, tokens chan token) {
	_, err := http.Get(url)
	if err != nil {
		// TODO: handle error
	}
	<-tokens
}

func main() {
	if len(os.Args) < 3 {
		exit()
	}

	url := os.Args[1]
	fmt.Println(url)
	limit, err := strconv.Atoi(os.Args[2])

	if err != nil {
		exit()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, os.Interrupt)

	tokens := make(chan token, limit)

	fmt.Println("Benchmark Start")
done:
	for {
		select {
		case <-stop:
			fmt.Println("Benchmark Stopping..")
			break done
		default:
			tokens <- token{}
			go work(url, tokens)
		}
	}

	// make sure existing goroutines have completed.
	for i := 0; i < limit; i++ {
		tokens <- token{}
	}

	close(stop)
	close(tokens)

	fmt.Println("Benchmark Done")
}
