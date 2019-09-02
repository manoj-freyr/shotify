// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"log"
	"sync"
	"utils"
	"encoding/json"
)

func main() {
	// simulating input urls
	urls := []string{"https://www.netflix.com","https://www.google.com","https://www.github.com", "https://www.medium.com"}
	var wg sync.WaitGroup
	log.Println("main enter")
	response_ch := make(chan utils.SSResponse)
	for _, url := range urls {
		wg.Add(1)
        go utils.Main_runner(url,&wg,response_ch)
    }
	results := make([]utils.SSResponse, len(urls))
	for i := range results {
        results[i] = <-response_ch
    }
	wg.Wait()
	jsonInfo, _ := json.Marshal(results)
	log.Printf("jsonInfo: %s\n", jsonInfo)
	return
}
