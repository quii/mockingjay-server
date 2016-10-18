package mockingjay

import "net/http"
import "math/rand"

var (
	httpMethods    = getHTTPMethods()
	urlLetterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randomURL(length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = urlLetterRunes[rand.Intn(len(urlLetterRunes))]
	}
	return string(b)
}

const maxURLLen = 2048 // can probably do more, but not sure value of it http://stackoverflow.com/questions/2659952/maximum-length-of-http-get-request

func getHTTPMethods() []string {
	return []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
}
