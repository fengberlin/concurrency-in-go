package main

import (
	"fmt"
	"math/rand"
	"time"
)

func simpleGoroutineLeak() {

	doWork := func(strs <-chan string) <-chan interface{} {
		completed := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(completed)
			for s := range strs {
				fmt.Println(s)
			}
		}()
		return completed
	}

	// 传递了一个 nil 进入 doWork, 导致strs永远无法读取到任何东西，
	// 而且包含doWork的goroutine将在这个过程的整个生命周期中保留在内存中
	doWork(nil)
	// some works here...
	fmt.Println("Done.")
}

func preventGoroutineLeak() {
	doWork := func(done <-chan interface{}, strs <-chan string) <-chan interface{} {
		terminated := make(chan interface{})
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(terminated)
			for {
				select {
				case s := <-strs:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()

		return terminated
	}

	done := make(chan interface{})
	terminated := doWork(done, nil)

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Canceling doWork goroutine...")
		close(done)
	}()

	<-terminated
	fmt.Println("Done.")
}

func goroutineNotStop() {
	newRandStream := func() <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.") // 1
			defer close(randStream)
			for {
				randStream <- rand.Int()
			}
		}()

		return randStream
	}

	randStream := newRandStream()
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
}

func stopGoroutine() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("newRandStream closure exited.")
			defer close(randStream)

			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()

		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}

	close(done)
	// do something
	time.Sleep(1 * time.Second)
}

func main() {
	// simpleGoroutineLeak()
	// preventGoroutineLeak()
	goroutineNotStop()
	// stopGoroutine()
}