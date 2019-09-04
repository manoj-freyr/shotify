package utils
import "sync"
// screen shot response from workers
type SSResponse struct{
	URL string
	Err error
	Data []byte
}

type SvcResponse struct {
  URL string `json:"url"`
  Status string `json:"status"`
  Link  string `json:"link"`
}

type LinkMap struct{
    // url---->link
    Lookup map[string]string
    CountersLock sync.RWMutex
}
func (mp *LinkMap) IsPresent(url string)(string, bool){
	mp.CountersLock.RLock()
	defer mp.CountersLock.RUnlock()
	if val, ok := mp.Lookup[url]; ok {
		return val,true
	}
	return "",false
}

func (mp *LinkMap) Insert(url,filename string){
    mp.CountersLock.Lock()
	defer mp.CountersLock.Unlock()
	mp.Lookup[url] = filename
}
