package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Result string
type Search func(query string) Result

var (
	Web   = fakeSearch("web")
	Image = fakeSearch("image")
	Video = fakeSearch("video")

	// All replicas
	WebReplica   = fakeSearch("web1")
	ImageReplica = fakeSearch("image1")
	VideoReplica = fakeSearch("video1")

	WebReplica2   = fakeSearch("web2")
	ImageReplica2 = fakeSearch("image2")
	VideoReplica2 = fakeSearch("video2")
)

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s result for %q\n", kind, query))
	}
}

// Query Web, Image, Video sequentially one after the other
func google1_0(query string) (results []Result) {
	results = append(results, Web(query))
	results = append(results, Image(query))
	results = append(results, Video(query))
	return
}

// Query Web, Image, and Video services concurrently
func google2_0(query string) (results []Result) {
	out := make(chan Result)
	searchMethod := []Search{Web, Video, Image}
	for _, v := range searchMethod {
		go func(s Search) {
			out <- s(query)
		}(v)
	}

	for range searchMethod {
		results = append(results, <-out)
	}

	return
}

// Query Web, Image, and Video services concurrently
// Services are given a timeout of 80ms
// If a service takes longer than 80ms, it is skipped
func google2_1(query string) (results []Result) {
	out := make(chan Result)
	searchMethod := []Search{Web, Video, Image}
	for _, v := range searchMethod {
		go func(s Search) {
			out <- s(query)
		}(v)
	}

	timeout := time.After(80 * time.Millisecond)
	for range searchMethod {
		select {
		case result := <-out:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Time out")
			return
		}
	}
	return
}

// Each service has multiple replicas
// All replicas are queried concurrently in addition to every service being queried concurrently
func google3_0(query string) (results []Result) {
	out := make(chan Result)
	searchMethods := [][]Search{
		{Web, WebReplica, WebReplica2},
		{Video, VideoReplica, VideoReplica2},
		{Image, ImageReplica, ImageReplica2},
	}

	for _, v := range searchMethods {
		go func(s ...Search) {
			out <- searchWithReplicas(query, s...)
		}(v...)
	}

	timeout := time.After(80 * time.Millisecond)
	for range searchMethods {
		select {
		case result := <-out:
			results = append(results, result)
		case <-timeout:
			fmt.Println("Time out")
			return
		}
	}
	return
}

// We get multiple replicas of the same Search (consider these different servers)
// We have a separate goroutine to search in each server
// While all the replicas search concurrently
// We simply return the first one that writes to the channel i.e. the fastest replica
func searchWithReplicas(query string, replicas ...Search) Result {
	out := make(chan Result)

	searchInReplica := func(replicaNum int) {
		out <- replicas[replicaNum](query)
	}

	for i := range replicas {
		go searchInReplica(i)
	}

	return <-out
}

func main() {
	start := time.Now()
	results := google1_0("golang")
	elapsed := time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	fmt.Print("---------------\n\n")

	start = time.Now()
	results = google2_0("golang")
	elapsed = time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	fmt.Print("---------------\n\n")

	start = time.Now()
	results = google2_1("golang")
	elapsed = time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)

	fmt.Print("---------------\n\n")

	start = time.Now()
	results = google3_0("golang")
	elapsed = time.Since(start)
	fmt.Println(results)
	fmt.Println(elapsed)
}
