// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"log"
	"sync"
	"utils"
	"encoding/json"
	"io/ioutil"
)

func main() {
	// simulating input urls
	urls := []string{"https://www.netflix.com","https://www.google.com","https://www.github.com", "https://www.medium.com"}
	var wg sync.WaitGroup
	log.Println("main enter")
	response_ch := make(chan *utils.SSResponse)
	for _, url := range urls {
		wg.Add(1)
        go utils.Main_runner(url,&wg,response_ch)
    }
	results := make([]utils.SvcResponse, len(urls))
	for i := range results {
		res := <-response_ch
		if res.Err != nil{
			results[i] = utils.SvcResponse{res.URL,res.Err.Error(),""}
		}else{
			link := utils.ConvertURL(res.URL)+".png"
			err := ioutil.WriteFile(link, res.Data, 0644)
			if err!= nil{
				results[i] = utils.SvcResponse{res.URL,err.Error(),""}
			}else{
				results[i] = utils.SvcResponse{res.URL,"success",link}
			}
		}
	}
	wg.Wait()
	jsonInfo, _ := json.Marshal(results)
	log.Printf("jsonInfo: %s\n", jsonInfo)
	return
}
