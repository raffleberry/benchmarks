package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type token struct{}

func exit() {
	fmt.Println("USAGE: bench.go <url> <conncurrency> <buffer-size>")
	os.Exit(1)
}

func work(url string, tokens <-chan token) {
	ok := true
	startTime := time.Now().UnixMilli()
	_, err := http.Get(url)
	if err != nil {
		ok = false
	}
	endTime := time.Now().UnixMilli()
	execTime := endTime - startTime
	fmt.Printf("%v,%v,%v\n", ok, endTime, execTime)
	<-tokens
}

func main() {
	if len(os.Args) < 3 {
		exit()
	}

	url := os.Args[1]
	limit, err1 := strconv.Atoi(os.Args[2])

	if err1 != nil {
		exit()
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, os.Interrupt)

	tokens := make(chan token, limit)

done:
	for {
		select {
		case <-stop:
			break done
		default:
			tokens <- token{}
			go work(url, tokens)
		}
	}

	for i := 0; i < limit; i++ {
		tokens <- token{}
	}
	close(stop)
	close(tokens)
}
