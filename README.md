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

