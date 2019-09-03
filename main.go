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
	archiveServer := http.NewServeMux()
	archiveServer.HandleFunc("/",utils.ArchiveHandler)
	go func() {
		http.ListenAndServe(":8008", archiveServer)
	}()
	http.ListenAndServe(":8002", screenShotServer)
}

