package main

import (
	"fmt"
	"time"
)

func whisper(left, right chan int) {
	rightVal := <-right
	fmt.Println("right: ", rightVal)

	time.Sleep(time.Second * 1)
	left <- rightVal + 1
	fmt.Println("left: ", rightVal+1)

}

/*
We start by creating a chain from leftmost to right.
As every item is added to the right, we call the `whisper`
goroutine, which blocks to receive a value at the channels.

Finally, we get to the point where we add a value to the
rightmost channel. This creates a chain that propagates the
value all the way to the leftmost channel.

Every left channel is a right channel to one other whisper().
So, once we add to the right most channel, it in turn adds to its left channel.
Which then adds to its left and so on and so forth. In this way, each left
channel blocks because it needs to receive input from its right.

The leftmost channel is not linked to another channel, that is why we block
it in the main function.
*/
func main() {
	const n = 15
	leftmost := make(chan int)
	left := leftmost
	right := leftmost
	for i := 0; i < n; i++ {
		right = make(chan int)
		go whisper(left, right)
		left = right
	}

	fmt.Println("Starting chain")
	go func(c chan int) {
		c <- 1
	}(right)

	// this makes sure that the leftmost channel blocks the main
	fmt.Println(<-leftmost)
}
