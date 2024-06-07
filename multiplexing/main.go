package main

import (
	"fmt"
)

func fanIn(input ...<-chan string) <-chan string {
	output := make(chan string)

	for _, ch := range input {
		go func(ch <-chan string) {
			for {
				val := <-ch
				output <- val
			}
		}(ch)
	}
	return output
}

func getDialogue(actor string) <-chan string {
	input := make(chan string)
	go func() {
		for i := 0; i < 10; i++ {
			input <- fmt.Sprintf("%v:  dialogue %v\n", actor, i)
		}
	}()

	return input
}

func main() {
	var channels []<-chan string

	for i := 0; i < 5; i++ {

		inp := getDialogue(fmt.Sprintf("Actor %v", i))
		channels = append(channels, inp)
	}

	result := fanIn(channels...)
	for i := 0; i < 5*10; i++ {
		fmt.Println(<-result)
	}
}
