package mockingjay

import (
	"bytes"
	"errors"
	"math/rand"
	"net/http"
	"strings"
)

const maxURLLen = 2048 // can probably do more, but not sure value of it http://stackoverflow.com/questions/2659952/maximum-length-of-http-get-request

var (
	httpMethods = getHTTPMethods()
	urlRunes    = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-_")
)

func randomURL(length uint16) (string, error) {
	pathLength := length / 2
	queryLength := length - pathLength
	path := randomPath(pathLength)
	query, err := randomQueryString(queryLength)
	if err != nil {
		return "", err
	}
	return path + query, nil
}

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

func randomPath(length uint16) (path string) {
	var p []rune
	p = append(p, '/')

	for len(p) < int(length) {
		for i := 0; i < rand.Intn(int(length)-len(p)); i++ {
			p = append(p, randomRune())
		}
		p = append(p, '/')
	}

	for len(p) < int(length) {
		p = append(p, randomRune())
	}

	return string(p)
}
func randomRune() rune {
	return urlRunes[rand.Intn(len(urlRunes))]
}

func randomQueryString(length uint16) (path string, err error) {
	if length < 5 {
		return path, errors.New("query string must be at least 5 characters long")
	}

	qsps := queryStringParameters{targetLength: int(length)}

	if length <= 10 {
		qsp := randomQueryStringParameter(int(length - 2))
		qsps.add(qsp)
		return qsps.join(), err
	}

	for int(length)-qsps.length() > 10 {
		qsp := randomQueryStringParameter(9)
		qsps.add(qsp)
	}

	qsps.padFinalQueryString()
	path = qsps.join()
	return
}

func randomQueryStringParameter(length int) (qsp queryStringParameter) {
	key := make([]rune, rand.Intn(length/2))
	value := make([]rune, length-len(key))
	fillWithRandomRunes(key)
	fillWithRandomRunes(value)

	return queryStringParameter{
		key:   string(key),
		value: string(value),
	}
}

func fillWithRandomRunes(slice []rune) {
	for i := range slice {
		slice[i] = randomRune()
	}
}

type queryStringParameter struct {
	key   string
	value string
}

func (qsp *queryStringParameter) join() string {
	return qsp.key + "=" + qsp.value
}

func (qsp *queryStringParameter) length() int {
	return len(qsp.join())
}

type queryStringParameters struct {
	targetLength int
	params       []queryStringParameter
}

func (qsps *queryStringParameters) add(qsp queryStringParameter) {
	qsps.params = append(qsps.params, qsp)
}

func (qsps *queryStringParameters) length() (length int) {
	if qsps.empty() {
		return 0
	}
	return len(qsps.join())
}

func (qsps *queryStringParameters) empty() bool {
	if len(qsps.params) == 0 {
		return true
	}
	return false
}

func (qsps *queryStringParameters) padFinalQueryString() {
	remaining := qsps.targetLength - qsps.length()
	additional := make([]rune, remaining)

	for i := range additional {
		additional[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	qsps.params[len(qsps.params)-1].key += string(additional)
}

func (qsps queryStringParameters) join() string {
	var b bytes.Buffer
	var queries []string
	b.WriteRune('?')
	for _, qsp := range qsps.params {
		queries = append(queries, qsp.join())
	}
	b.WriteString(strings.Join(queries, "&"))
	return b.String()
}
