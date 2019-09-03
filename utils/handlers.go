package utils
import(
	"net/http"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"bufio"
	"sync"
	"log"
	"encoding/json"
	"strings"
)

func archiveHandler(w http.ResponseWriter, r *http.Request){
	b, _ := ioutil.ReadAll(r.Body)
	fileName,_ := r.URL.Query()["fileName"]
	outer,inner := GetFolders(fileName[0])
	os.MkdirAll(outer+inner, 0644)
	_ = ioutil.WriteFile(outer+inner+fileName[0], b, 0644)
}

//for GETs
func queryHandler(w http.ResponseWriter, r *http.Request){
    fileName,_ := r.URL.Query()["fileName"]
    outer,inner := GetFolders(fileName[0])
	http.ServeFile(w,r,outer+inner+fileName[0])
}

func errorHandler(w http.ResponseWriter, r *http.Request, err string){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
	return
}

func readFromFile(r *http.Request) ([]string,error) {
	b, err := ioutil.ReadAll(r.Body)
	if err!= nil{
		fmt.Println("error in reading")
        return nil,err
     }
     rd := bytes.NewReader(b)
	 urls := []string{}
     scanner := bufio.NewScanner(rd)
     scanner.Split(bufio.ScanLines)
     for scanner.Scan() {
         urls = append(urls, scanner.Text())
     }
	return urls,nil
}

func ListHandler(w http.ResponseWriter, r *http.Request){
    if r.Method == "GET" {
		errorHandler(w,r,"400 Get not supported!")
        return
    }
	var urls []string
    urlslist, listok := r.URL.Query()["urls"]
	_, fileok := r.URL.Query()["urlfile"]
	if (listok && fileok) || (!listok && !fileok){
		errorHandler(w,r,"400 - Either specify url list OR use file!")
		return
	}
	if (listok && len(urls)<1){
		errorHandler(w,r,"400 - No Urls specified!") 
        return
    }
	if listok{
		urls = strings.Split(urlslist[0] ,",")
	}
	var err error
	if fileok{
		urls,err = readFromFile(r)
		if err!=nil{
			errorHandler(w,r,"400 - Bad url file")
			return
		}
	}
    fmt.Println("urls are", urls)
    //urls := []string{"https://www.netflix.com","https://www.google.com","https://www.github.com", "https://www.medium.com"}
    var wg sync.WaitGroup
    log.Println("main enter")
    response_ch := make(chan *SSResponse)
    for _, url := range urls {
        wg.Add(1)
        go Main_runner(url,&wg,response_ch)
    }
    results := make([]SvcResponse, len(urls))
    for i := range results {
        res := <-response_ch
        if res.Err != nil{
            results[i] = SvcResponse{res.URL,res.Err.Error(),""}
        }else{
            link := ConvertURL(res.URL)+".png"
            err := ioutil.WriteFile(link, res.Data, 0644)
            //err := uploadToArchive(res.Data,link,"127.0.0.1:8008")
            //err := http.NewRequest("POST", "127.0.0.1:8008", res.Data)
            if err!= nil{
                results[i] = SvcResponse{res.URL,err.Error(),""}
            }else{
                results[i] = SvcResponse{res.URL,"success",link}
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
