package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	start := time.Now()
	urls := os.Args[1:]
	ch := make(chan string)
	for _, url := range urls {
		go fetch(url, ch)
	}
	for range os.Args[1:] {
		fmt.Println(<-ch)
	}
	fmt.Printf("total time: %.4fs\n", time.Since(start).Seconds())
}
func fetch(url string, ch chan<- string) {
	start := time.Now()
	res, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("failed to fetch %s, %v", url, err)
	}
	defer res.Body.Close()
	name := strconv.Itoa(time.Now().Nanosecond())
	f, err := os.Create(name)
	if err != nil {
		ch <- fmt.Sprintf("failed to create file, %s, %v", url, err)
	}
	defer f.Close()
	size, err := io.Copy(f, res.Body)
	if err != nil {
		ch <- fmt.Sprintf("failed to read response body, %s, %v", url, err)
	}
	ch <- fmt.Sprintf("%d, %.4fs", size, time.Since(start).Seconds())
}
