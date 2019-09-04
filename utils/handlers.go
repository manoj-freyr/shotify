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

var lookUpMap = LinkMap{
    Lookup: make(map[string]string),
}


func ArchiveHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		queryHandler(w,r)
	}else{
		uploadHandler(w,r)
	}
}
func uploadHandler(w http.ResponseWriter, r *http.Request){
	b, _ := ioutil.ReadAll(r.Body)
	fileName,_ := r.URL.Query()["fileName"]
	link := ConvertURL(fileName[0])
	outer,inner := GetFolders(link)
	os.MkdirAll(outer+"/"+inner, 0755)
	err := ioutil.WriteFile(outer+"/"+inner+"/"+link+".png", b, 0644)
	if err!= nil{
		fmt.Println(err)
	}
	log.Println("in uploadHandler"+fileName[0])
	lookUpMap.Insert(fileName[0], link)
}

//for GETs
func queryHandler(w http.ResponseWriter, r *http.Request){
    fileName,err := r.URL.Query()["fileName"]
	if err != true{
		errorHandler(w,r,"Invalid URL, no fileName")
		return
	}
	link := ConvertURL(fileName[0])
    outer,inner := GetFolders(link)
	http.ServeFile(w,r,outer+"/"+inner+"/"+link+".png")
}

func errorHandler(w http.ResponseWriter, r *http.Request, err string){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err))
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


func writeHelper(url,filename string, data []byte)(string, error){
	resp,err := http.Post("http://127.0.0.1:8008/req?fileName="+filename,"application/octet-stream" ,bytes.NewReader(data))
	if err!=nil{
		log.Println("who am i")
		return "",err
	}
	defer resp.Body.Close()
	outer,inner := GetFolders(filename)
	return outer+inner+filename, nil
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
	if (listok && len(urlslist)<1){
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
	var cachedresults []SvcResponse
    fmt.Println("urls are", urls)
    var wg sync.WaitGroup
    log.Println("main enter")
    response_ch := make(chan *SSResponse)

	newCnt := 0
    for _, url := range urls {
		log.Println("url is"+url)
		imLink,ok := lookUpMap.IsPresent(url)
		if !ok{
			newCnt++
			wg.Add(1)
			go Main_runner(url,&wg,response_ch)
		}else{
			fmt.Println("Content already in")
			cachedresults = append(cachedresults,SvcResponse{url,"success",imLink})
		}
    }
    results := make([]SvcResponse, newCnt)
    for i := range results {
        res := <-response_ch
        if res.Err != nil{
            results[i] = SvcResponse{res.URL,res.Err.Error(),""}
        }else{
            link := ConvertURL(res.URL)
			_,err:= writeHelper("127.0.0.1:8008",link,res.Data)
            if err!= nil{
                results[i] = SvcResponse{res.URL,err.Error(),""}
            }else{
                results[i] = SvcResponse{res.URL,"success",link}
            }
        }
    }
    wg.Wait()
	finalResults := append(results,cachedresults...)
	jsonInfo, _ := json.Marshal(finalResults)
    log.Printf("jsonInfo: %s\n", jsonInfo)
    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonInfo)
    return
}
