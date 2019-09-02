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
