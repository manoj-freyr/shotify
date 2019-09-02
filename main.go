// Command screenshot is a chromedp example demonstrating how to take a
// screenshot of a specific element and of the entire browser viewport.
package main

import (
	"log"
	"sync"
	"utils"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"fmt"
)

//var instance *LinkMap
func main(){
	screenShotServer := http.NewServeMux()
	screenShotServer.HandleFunc("/urls", listHandler)
	//screenShotServer.HandleFunc("/urlfile", fileHandler)
	//archiveServer := http.NewServeMux()
	//archiveServer.HandleFunc("/query",archiveHandler)
	//go func() {
	//	http.ListenAndServe(":8008", archiveServer)
	//}()
	http.ListenAndServe(":8002", screenShotServer)
}

func archiveHandler(w http.ResponseWriter, r *http.Request){
	u, err := r.URL.Parse(s)
	err := ioutil.WriteFile(u.Path, res.Data, 0644)
	
}

func uploadToArchive(data []byte, fileName , url string) error{
	fullName := url+"?name="+fileName
	request, err := http.NewRequest("POST", fullName, data)
	client := &http.Client{}
	request.Header.Add("Content-Type", "image/png")
	response, err := client.Do(request)
	defer response.Body.Close()
	return nil

}
func listHandler(w http.ResponseWriter, r *http.Request){
	urls, listok := r.URL.Query()["urls"]
	if !listok || len(urls)<1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("500 - No Urls specified!"))
		return
	}
	fmt.Println("urls are", urls)
	//urls := []string{"https://www.netflix.com","https://www.google.com","https://www.github.com", "https://www.medium.com"}
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
			//err := uploadToArchive(res.Data,link,"127.0.0.1:8008")
			//err := http.NewRequest("POST", "127.0.0.1:8008", res.Data)
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonInfo)
	return
}
