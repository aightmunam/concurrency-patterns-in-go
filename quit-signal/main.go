package main

import (
	"fmt"
	"math/rand"
	"time"
)

func getDialogue(quit chan bool) chan string {
	output := make(chan string)
	go func(ch chan string) {
		for i := 0; ; i++ {
			select {
			case output <- fmt.Sprintf("Dialogue # %v", i):
			case <-quit:
				duration := time.Duration(rand.Intn(1e3)) * time.Millisecond
				time.Sleep(duration)
				fmt.Printf("Did some clean up before quitting that took: %v\n", duration)
				quit <- true
				return
			}
		}

	}(output)
	return output
}

func main() {

	// two way communication to quit
	// `main` tells the `getDialogue` to quit
	// `getDialogue` gets the quit message and:
	//	1. Does all the cleanup related to quitting
	// 	2. Lets the `main` know again right before quitting
	//	3. Actually quits
	quit := make(chan bool)
	ch := getDialogue(quit)

	for i := 3; i >= 0; i-- {
		fmt.Println(<-ch)
	}
	quit <- true
	fmt.Println("Did you quit", <-quit)
}
