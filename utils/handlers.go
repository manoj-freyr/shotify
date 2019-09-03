package utils
import(
	"http"
	"fmt"
)

func errorHandler(w http.ResponseWriter, r *http.Request, err string){
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(err)
	return
}

func readFromFile(r *http.Request) ([]string,error) {
	b, err := ioutil.ReadAll(r.Body)
	if err!= nil{
		fmt.Println("error in reading")
        return nil,err
     }
     r := bytes.NewReader(b)
	 urls := []string{}
     scanner := bufio.NewScanner()
     scanner.Split(bufio.ScanLines)
     for scanner.Scan() {
         urls.append(result, scanner.Text())
     }
	return urls,nil
}

func ListHandler(w http.ResponseWriter, r *http.Request){
    if r.Method == "GET" {
		errorHandler(w,r,"400 Get not supported!")
        return
    }
    urls, listok := r.URL.Query()["urls"]
	file, fileok := r.URL.Query()["urlfile"]
	if (listok && fileok) || (!listok && !fileok){
		errorHandler(w,r,"400 - Either specify url list OR use file!")
		return
	}
	if (listok && len(urls)<1){
		errorHandler(w,r,"400 - No Urls specified!") 
        return
    }
	if listok{
		urls = strings.Split(urls ",")
	}
	if fileok{
		urls,err := readFromFile(r)
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
