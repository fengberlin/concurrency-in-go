package main

import (
	"bytes"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	// 条件等待
	cadence := sync.NewCond(&sync.Mutex{})

	go func() {
		for range time.Tick(1 * time.Millisecond) {
			cadence.Broadcast()
		}
	}()

	takeStep := func() {
		cadence.L.Lock()
		cadence.Wait()
		cadence.L.Unlock()
	}

	tryDirection := func(direction string, moveCount *int32, out *bytes.Buffer) bool {
		fmt.Fprintf(out, " %v", direction)
		atomic.AddInt32(moveCount, 1)
		takeStep()
		if atomic.LoadInt32(moveCount) == 1 {
			fmt.Fprint(out, ". Success!")
			return true
		}
		takeStep()
		atomic.AddInt32(moveCount, -1)
		return false
	}

	var left, right int32
	tryLeft := func(out *bytes.Buffer) bool {
		return tryDirection("left", &left, out)
	}
	tryRight := func(out *bytes.Buffer) bool {
		return tryDirection("right", &right, out)
	}

	walk := func(walking *sync.WaitGroup, name string) {
		var out bytes.Buffer
		defer func() {
			fmt.Println(out.String())
		}()
		defer walking.Done()
		// scoot: 溜走
		fmt.Fprintf(&out, "%v is trying to scoot:", name)
		for i := 0; i < 5; i++ {
			if tryLeft(&out) || tryRight(&out) {
				return
			}
		}
		fmt.Fprintf(&out, "\n%v tosses her hands up in exasperation!", name)
	}

	var peopleInHallway sync.WaitGroup
	(&peopleInHallway).Add(2)
	go walk(&peopleInHallway, "Alice")
	go walk(&peopleInHallway, "Bob")
	(&peopleInHallway).Wait()
}
