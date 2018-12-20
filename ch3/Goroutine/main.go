package main

import "fmt"

func sayHello(done chan<- struct{}) {
	fmt.Println("hello")
	done <- struct{}{}
}

func main() {
	done := make(chan struct{})
	go sayHello(done)
	<-done
}