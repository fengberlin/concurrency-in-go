package main

import (
	"fmt"
	"sync"
)

func main() {
	var count int
	increment := func() {
		count++
	}

	var once sync.Once
	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			// sync.Once 只计算 Do 被调用的次数
			once.Do(increment)
		}()
	}

	wg.Wait()
	fmt.Printf("Count is %d\n", count)
}