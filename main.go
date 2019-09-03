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
	screenShotServer.HandleFunc("/", utils.ListHandler)
	screenShotServer.HandleFunc("/urlFile", fileHandler)
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
