package mockingjay

import (
	"math/rand"
	"net/http"
)

var (
	httpMethods = getHTTPMethods()
	urlRunes    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
)

func randomURL(length uint8) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = urlRunes[rand.Intn(len(urlRunes))]
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

func randomPath(length uint8) (path string) {
	var p []rune

	for len(p) < int(length) {
		for i := 0; i < rand.Intn(int(length)-len(p)); i++ {
			p = append(p, urlRunes[rand.Intn(len(urlRunes))])
		}
		p = append(p, '/')
	}

	for len(p) < int(length) {
		p = append(p, urlRunes[rand.Intn(len(urlRunes))])
	}

	return string(p)
}

func randomQueryString(length uint8) (path string) {
	var p []rune
	p = append(p, '?')

	for i := 0; i < rand.Intn(int(length/4)+1); i++ {
		p = append(p, urlRunes[rand.Intn(len(urlRunes))])
	}
	p = append(p, '=')
	for i := 0; i < rand.Intn(int(length/4)+1); i++ {
		p = append(p, urlRunes[rand.Intn(len(urlRunes))])
	}

	for len(p) < int(length) {
		if len(p)+2 >= int(length) {
			return string(p)
		}
		p = append(p, '&')
		for i := 0; i < rand.Intn((int(length)-len(p))/2); i++ {
			p = append(p, urlRunes[rand.Intn(len(urlRunes))])
		}
		p = append(p, '=')
		for i := 0; i < rand.Intn(int(length)-len(p)); i++ {
			p = append(p, urlRunes[rand.Intn(len(urlRunes))])
		}
	}

	for len(p) < int(length) {
		p = append(p, urlRunes[rand.Intn(len(urlRunes))])
	}

	return string(p)
}
