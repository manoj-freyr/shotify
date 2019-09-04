package utils
import(
	"strings"
	"hash/fnv"
    "strconv"
	 "net/url"
)

// a small hack, as chromedp wont accept without http or https in url
func HTTPify(urlstr string)(string,bool){
	u,err := url.Parse(urlstr)
	appendstr:= ""
	if err!= nil{
		return "",false
	} else if (u.Host == "") ||(u.Scheme == "") {
		if strings.Index(urlstr,"//") != -1{
			appendstr = "http:"
		}else{
			appendstr = "https://"
		}
	}
	return appendstr+urlstr,true
}

func ConvertURL(url string) string{
    // Create replacer with pairs as arguments.
    r := strings.NewReplacer("/", "%2F",
            ":", "%3A")

    return r.Replace(url)
}


func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func GetFolders(url string) (string,string) {
	hashNum := hash(url)
	return strconv.Itoa(int(hashNum) % 100), strconv.Itoa(int(hashNum) % 1000)
}

