package main

import (
	"fmt"
)

type Message struct {
	text string
	wait chan bool
}

func fanIn(input ...<-chan Message) <-chan Message {
	output := make(chan Message)

	for _, c := range input {
		go func(c <-chan Message) {
			for {
				val := <-c
				output <- val
			}
		}(c)
	}
	return output
}

func getDialogue(actor string) <-chan Message {
	input := make(chan Message)
	wait := make(chan bool)

	go func() {
		for i := 0; i < 10; i++ {
			message := Message{
				text: fmt.Sprintf("%v:  dialogue %v\n", actor, i),
				wait: wait,
			}
			input <- message

			// this will block an actor after they have spoken their dialogue
			// the actor cannot move on to the next dialogue till the wait
			// channel receives some message
			<-message.wait
		}
	}()

	return input
}

func main() {
	var channels []<-chan Message
	const totalActors = 3

	for i := 0; i < totalActors; i++ {
		inp := getDialogue(fmt.Sprintf("Actor %v", i))
		channels = append(channels, inp)
	}

	c := fanIn(channels...)

	// Make sure all the actors have spoken their dialogue before we move to the
	// next dialogue

	for i := 0; i < 10; i++ {
		messages := make([]Message, totalActors)

		// Get the dialogue for every actor
		// Every actor gets blocked by their wait channel after their dialogue
		for j := 0; j < totalActors; j++ {
			messages[j] = <-c
			fmt.Println(messages[j].text)
		}

		// Unblock all actors by passing a message to their wait channel
		for j := 0; j < 3; j++ {
			messages[j].wait <- false
		}
		fmt.Println("-----------------")
	}
}
