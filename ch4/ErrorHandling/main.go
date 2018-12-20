package main

import (
	"fmt"
	"net/http"
)

// 仅仅是展示了错误出来，而没有选择做更多的处理
func limitedHandling() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan *http.Response {
		responses := make(chan *http.Response)
		go func() {
			defer close(responses)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case responses <- resp:
				}
			}
		}()

		return responses
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"https://www.baidu.com", "https://badhost"}
	for response := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", response.Status)
	}
}

type Result struct {
	Error    error
	Response *http.Response
}

func improvedHandling() {
	checkStatus := func(done <-chan interface{}, urls ...string) <-chan Result {

		results := make(chan Result)
		go func() {
			defer close(results)

			for _, url := range urls {
				var result Result
				resp, err := http.Get(url)
				result = Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case results <- result:
				}
			}
		}()

		return results
	}

	done := make(chan interface{})
	defer close(done)

	// urls := []string{"https://www.baidu.com", "https://badhost"}
	// for result := range checkStatus(done, urls...) {
	// 	if result.Error != nil {
	// 		fmt.Printf("error: %v\n", result.Error)
	// 		continue
	// 	}
	// 	fmt.Printf("Response: %v\n", result.Response.Status)
	// }

	errCount := 0
	urls := []string{"cc", "uuuu", "https://www.baidu.com", "https://badhost", "a", "bb"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error: %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("Too many errors, breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}
}

func main() {
	// limitedHandling()
	improvedHandling()
}
