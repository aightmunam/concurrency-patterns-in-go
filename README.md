## My takeaways:

The following are my takeaways in no particular order. 
Warning: Some of this might be utter nonsense except that
it helps me understand some of the concepts.


---
#### Generator Pattern:
A pattern where a function returns a send channel i.e. channel that can only send data out but cannot receive data.

Such a function usually calls another goroutine inside it (can be an anonymous function goroutine) and that goroutine communicates with the channel.

Example code below:
```
func generator() (<-chan string) {
    outputCh := make(chan string)
    go func() {
        outputCh <- "Some random shit"
    }()

    return outputCh
}
```

---
#### Asleep Goroutines 

If a channel is waiting to receive some data, but all the goroutines
are closed i.e. no body is sending data anymore, this will result in a
`fatal error: all goroutines are asleep - deadlock!` error.
Let's try to reproduce this:
```
// this function returns a channel that will return 5 messages only
func dataProvider() (<-chan string) {
    datach := make(chan string)
    go func() {
        for i := 0; i < 5; i++ {
            datach <- fmt.Sprintf("Message # %v", i)
        }
    }()
    return datach
}

func main() {
    datach := dataProvider()

    // this channel is expecting 10 values
    for i := 0; i < 10; i++ {
        fmt.Println(<-datach)
    }
}
```
In such a situation, the goroutine will close after sending 5 messages
to the channel, but the channel is expecting 10 messages. When the 
goroutine closes, the compiler notices all the goroutines are dead but the
channel is still open and expecting more messages so it raises a fatal error.

---

#### Data left in the channel
What happens if the channel has more data left in it but it is never received or read from it?

Example below:
```
// Returns a channel that will have three messages added to it
func dataGetter() (<-chan string) {
    ch := make(chan string)
    go func() {
        ch <- "Message 1"
        ch <- "Message 2"
        ch <- "Message 3"
    }()
    return ch
}

// Reads a single message from the chan
func main() {
    ch := dataGetter()
    fmt.Println(<-ch)
}

```
In the above example, *no error is shown*. The main function simply reads one message from the channel and closes. The other two messages are ignored.


---

### Some quotes directly from [The Go Memory Model](https://go.dev/ref/mem)

#### Initialization

>If a package p imports package q, the completion of q's init functions happens before the start of any of p's.

Go waits for the completion of all init functions for all the imported packages before the main function is called.

#### Goroutine

>The go statement that starts a new goroutine is synchronized before the start of the goroutine's execution. However, the exit of a goroutine is not guaranteed to be synchronized before any event in the program.

Creating a goroutine happens before the goroutine’s execution begins. There is no such guarantee incase of exit of a goroutine.

Take the following example:
```go
i := 0
go func() {
    i++
}()
fmt.Println(i)
```
No race condition between `i := 0` and `i++` since the goroutine needs to be created before it is executed. And it is created after `i` is initialized.

Race condition between `i++` and `fmt.Println(i)` because we have no way knowing which of the two statements will be called first. 

#### Channels

> A send on a channel is synchronized before the completion of the corresponding receive from that channel.

A send on a channel happens before the receive on the channel.
```go
i := 0
ch := make(chan int)
go func() {
    <-ch
    fmt.Println(i)
}()

i++
ch <- 0
```
The order is as follows:
*variable increment < channel send < channel receive < variable read*


<br>


>The closing of a channel is synchronized before a receive that returns a zero value because the channel is closed.

```go
var c = make(chan int, 10)
var a string

func f() {
	a = "hello, world"
	c <- 0

    // closing a channel sends the zero value to the channel
    // In this case, close(c) will send 0 (zero value for int) to the channel
    // c <- 0 is equivalent to close(c) in this case
}


func main() {
	go f()
	<-c  // synchronizes, so guarantees that hello, world will be printed

	print(a)
}
```

In the previous example, replacing `c <- 0` with `close(c)` yields a program with the same guaranteed behavior.

<br>

>A receive from an unbuffered channel is synchronized before the completion of the corresponding send on that channel.

A receive from an unbuffered channel happens before the send on that channel completes.

```go
var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c // will wait till the channel receives some message
}

func main() {
	go f()
	c <- 0
	print(a)  // guaranteed to print "hello, world" because the receive of a channel makes it wait till it receives some message
}
```

<br>

> The kth receive on a channel with capacity C is synchronized before the completion of the k+Cth send from that channel completes.

A channel with buffer size C, will receive its kth message before the C+kth message send completes.

- A counting semaphore can be modeled by a buffered channel.
  - the number of items in the channel corresponds to the number of active uses.
  - the capacity of the channel corresponds to the maximum number of simultaneous uses.
  - sending an item acquires the semaphore.
  - receiving an item releases the semaphore.
- A counting semaphore can be used to limit concurrency.

```go
var limit = make(chan int, 3) // buffered channel can be used to block after the buffer is full

func main() {
	for _, w := range work {
		go func(w func()) {
			limit <- 1 // each goroutine fills the buffer by 1. This blocks after 3 messages have been sent. 
			w()
			<-limit  // Each time we receive the message, the buffer gets a single capacity back
		}(w)
	}
	select{}
}
```

---