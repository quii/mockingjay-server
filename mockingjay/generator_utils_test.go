package mockingjay

import (
	"net/url"
	"testing"
	"testing/quick"
)

func TestRandomURL(t *testing.T) {
	f := func(n uint16) bool {
		path, _ := randomURL(n)
		_, err := url.Parse(path)
		return err == nil
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomUrl did not return a valid URL")
	}
}

func TestRandomPath(t *testing.T) {
	f := func(n uint16) bool {
		n += 1
		path := randomPath(n)
		_, err := url.Parse(path)

		if err != nil {
			t.Log(err)
			return false
		}

		if path[:1] != "/" {
			t.Log("First character of path is not /")
			return false
		}

		if len(path) != int(n) {
			t.Log("Wrong length of path returned")
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRandomQueryString(t *testing.T) {
	f := func(n uint16) bool {
		n = (n % 101) + 5

		qs, err := randomQueryString(n)
		if err != nil {
			t.Log(err)
			return true
		}

		if len(qs) != int(n) {
			t.Log("Wrong length of query string returned")
			return false
		}

		_, err = url.Parse(qs)

		if err != nil {
			t.Log(err)
			return false
		}

		if qs[:1] != "?" {
			t.Log("Query string doesn't start with a ?")
			return false
		}

		if !contains('=', qs) {
			t.Log("Query string has no =")
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestRandomQueryStringError(t *testing.T) {
	f := func(n uint16) bool {
		n %= 5
		_, err := randomQueryString(n)

		if err == nil {
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomQueryString did not produce an error with invalid input")
	}
}

func TestQueryStringParameter(t *testing.T) {
	qs := queryStringParameter{"hello", "world"}

	if qs.length() != 11 {
		t.Error("Unexpected length ", qs.length())
	}

	if qs.join() != "hello=world" {
		t.Error("Unexpected join value", qs.join())
	}
}

func TestQueryStringParameters(t *testing.T) {
	qsp := queryStringParameters{}
	if qsp.length() != 0 {
		t.Error("Unexpected length", qsp.length())
	}

	qsp.add(queryStringParameter{"hello", "world"})
	qsp.add(queryStringParameter{"goodbye", "piccadilly"})

	if qsp.length() != 31 {
		t.Error("Unexpected length", qsp.length())
	}

	qsp.add(queryStringParameter{"a", "b"})

	if qsp.join() != "?hello=world&goodbye=piccadilly&a=b" {
		t.Error("Unexpected join value", qsp.join())
	}
}

func contains(char rune, s string) bool {
	for _, r := range s {
		if r == char {
			return true
		}
	}
	return false
}
