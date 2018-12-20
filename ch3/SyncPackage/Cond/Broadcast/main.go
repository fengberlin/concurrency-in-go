package main

import (
	"fmt"
	"sync"
)

type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var tempwg sync.WaitGroup
		tempwg.Add(1)
		go func() {
			tempwg.Done()
			c.L.Lock()
			defer c.L.Unlock()
			c.Wait()
			fn()
		}()
		tempwg.Wait()
	}

	var wg sync.WaitGroup
	wg.Add(3)
	subscribe(button.Clicked, func() { //4
		fmt.Println("Maximizing window.")
		defer wg.Done()
	})
	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		defer wg.Done()
	})
	subscribe(button.Clicked, func() { //6
		fmt.Println("Mouse clicked.")
		defer wg.Done()
	})

	button.Clicked.Broadcast()

	wg.Wait()
}
