package utils
import(
	"strings"
)

func convertURL(url string) string{
    // Create replacer with pairs as arguments.
    r := strings.NewReplacer("/", "%2F",
            ".", "%2E",
            ":", "%3A")

    return r.Replace(url)
}
