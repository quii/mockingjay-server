package mockingjay

import (
	"bytes"
	"errors"
	"math/rand"
	"net/http"
	"strings"
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

type queryStringParameter struct {
	key   string
	value string
}

type queryStringParameters struct {
	targetLength int
	params       []queryStringParameter
}

func (qsp *queryStringParameter) join() string {
	return qsp.key + "=" + qsp.value
}

func (qsp *queryStringParameter) length() int {
	return len(qsp.join())
}

func (qsps *queryStringParameters) add(qsp queryStringParameter) {
	qsps.params = append(qsps.params, qsp)
}

func (qsps *queryStringParameters) addParamIfFit(qsp queryStringParameter) error {
	if qsp.length()+qsps.length() > qsps.targetLength {
		return errors.New("parameter longer than target length of query string")
	}
	qsps.add(qsp)
	return nil
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

func (qsps *queryStringParameters) full() bool {
	return len(qsps.params) >= qsps.targetLength
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
	var joinedStrings []string
	b.WriteRune('?')
	for _, qsp := range qsps.params {
		joinedStrings = append(joinedStrings, qsp.join())
	}
	b.WriteString(strings.Join(joinedStrings, "&"))
	return b.String()
}

func randomQueryString(length uint8) (path string, err error) {
	qsps := queryStringParameters{targetLength: int(length)}

	if length <= 10 {
		qsp, _ := randomQueryStringParameter(int(length - 2))
		qsps.add(qsp)
		return qsps.join(), err
	}

	for int(length)-qsps.length() > 10 {
		qsp, _ := randomQueryStringParameter(9)
		qsps.add(qsp)
	}

	qsps.padFinalQueryString()

	return qsps.join(), err
}

func randomQueryStringParameter(length int) (qsp queryStringParameter, err error) {
	if length < 3 {
		return qsp, errors.New("Minimum length is 3")
	}
	key := make([]rune, rand.Intn(length/2)+1)
	value := make([]rune, length-len(key))

	for i := range key {
		key[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	for i := range value {
		value[i] = urlRunes[rand.Intn(len(urlRunes))]
	}

	return queryStringParameter{string(key), string(value)}, nil
}
