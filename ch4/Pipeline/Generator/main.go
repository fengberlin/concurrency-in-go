package main

import (
	"fmt"
	"math/rand"
)

func repeatGenerator() {

	done := make(chan interface{})
	defer close(done)

	for num := range takeFunc(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
	fmt.Println()
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case valueStream <- v:
				}
			}
		}
	}()
	return valueStream
}

func takeFunc(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {

	takeStream := make(chan interface{})
	go func() {
		defer close(takeStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case takeStream <- <-valueStream:
			}
		}
	}()
	return takeStream
}

func repeatFnGenerator() {

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} {
		return rand.Int()
	}

	for num := range takeFunc(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}

func repeatFn(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func toStringFunc(done <-chan interface{}, valueStream <-chan interface{}) <-chan string {

	stringStream := make(chan string)
	go func() {
		defer close(stringStream)
		for v := range valueStream {
			select {
			case <-done:
				return
			case stringStream <- v.(string):
			}
		}
	}()
	return stringStream
}

func toStringGenerator() {
	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toStringFunc(done, takeFunc(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...\n", message)
}

func main() {
	// repeatGenerator()
	// repeatFnGenerator()
	toStringGenerator()
}
