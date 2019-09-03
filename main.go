package main

import (
	"utils"
	"net/http"
)

//var instance *LinkMap
func main(){
	screenShotServer := http.NewServeMux()
	screenShotServer.HandleFunc("/", utils.ListHandler)
	//screenShotServer.HandleFunc("/urlfile", fileHandler)
	//archiveServer := http.NewServeMux()
	//archiveServer.HandleFunc("/query",archiveHandler)
	//go func() {
	//	http.ListenAndServe(":8008", archiveServer)
	//}()
	http.ListenAndServe(":8002", screenShotServer)
}

/*
func uploadToArchive(data []byte, fileName , url string) error{
	fullName := url+"?name="+fileName
	request, err := http.NewRequest("POST", fullName, data)
	client := &http.Client{}
	request.Header.Add("Content-Type", "image/png")
	response, err := client.Do(request)
	defer response.Body.Close()
	return nil

}
*/
