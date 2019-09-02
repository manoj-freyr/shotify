package utils
 
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
