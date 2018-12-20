package main

import (
	"bytes"
	"fmt"
	"os"
)

func testingChannel() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}

func channelOwner() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		fmt.Println("writing data...")
		go func() {
			defer close(resultStream)
			defer fmt.Println("writing data is done.")
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving!")
}

func main() {
	// testingChannel()
	channelOwner()
}