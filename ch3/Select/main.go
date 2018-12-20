package main

import (
	"fmt"
	"time"
)

func selectExample() {
	start := time.Now()
	c := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(c)
	}()

	fmt.Println("Blocking on read...")
	select {
	case <-c:
		fmt.Printf("Unblocked %v later.\n", time.Since(start))
	}
}

func randomSelect() {
	c1 := make(chan interface{})
	close(c1)
	c2 := make(chan interface{})
	close(c2)

	var c1Count, c2Count int
	for i := 0; i < 1000; i++ {
		select {
		case <-c1:
			c1Count++
		case <-c2:
			c2Count++
		}
	}

	fmt.Printf("c1Count: %d\nc2Count: %d\n", c1Count, c2Count)
}

func timeoutSelect() {
	var c <-chan int
	select {
	case <-c: // 永久阻塞，因为从一个未初始化即为nil的channel读取
	case <-time.After(1 * time.Second):
		fmt.Println("Timed out.")
	}
}

func selectWithDefault() {
	start := time.Now()
	var c1, c2 <-chan int
	select {
	case <-c1:
	case <-c2:
	default:
		fmt.Printf("In default after %v\n", time.Since(start))
	}
}

func loopSelect() {
	done := make(chan interface{})
	go func() {
		time.Sleep(5 * time.Second)
		close(done)
	}()

	workCounter := 0
	loop:
	for {
		select {
		case <-done:
			break loop
		default:
			fmt.Println("default...")
		}

		workCounter++
		time.Sleep(1 * time.Second)
	}

	fmt.Printf("Achieved %v cycles of work before signalled to stop.\n", workCounter)
}

func main() {
	// selectExample()
	// randomSelect()
	// timeoutSelect()
	// selectWithDefault()
	loopSelect()
}
