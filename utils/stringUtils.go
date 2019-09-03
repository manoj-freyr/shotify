package utils
import(
	"strings"
	"hash/fnv"
    "strconv"
)

func ConvertURL(url string) string{
    // Create replacer with pairs as arguments.
    r := strings.NewReplacer("/", "%2F",
            ".", "%2E",
            ":", "%3A")

    return r.Replace(url)
}


func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func GetFolders(url string) (string,string) {
        hashNum := hash("HelloWorld")
        fmt.Println(strconv.Itoa(int(hashNum % 1000)), strconv.Itoa(int(hashNum % 100)))
}

