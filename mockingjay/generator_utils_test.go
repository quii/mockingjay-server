package mockingjay

import (
	"net/url"
	"testing"
	"testing/quick"
)

func TestRandomURL(t *testing.T) {
	f := func(n uint8) bool {
		path := randomURL(n)
		_, err := url.Parse(path)
		return err == nil
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomUrl did not return a valid URL")
	}
}

func TestRandomPath(t *testing.T) {
	f := func(n uint8) bool {
		path := randomPath(n)
		_, err := url.Parse(path)
		hasSlash := contains('/', path)

		return err == nil &&
			hasSlash &&
			len(path) == int(n)
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomPath did not did not return a valid path")
	}
}

func TestRandomQueryString(t *testing.T) {
	f := func(n uint8) bool {
		n = n + 2
		qs := randomQueryString(n)
		_, err := url.Parse(qs)

		if err != nil {
			return false
		}

		if qs[:1] != "?" {
			return false
		}

		if !contains('=', qs) {
			return false
		}

		return true
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomQueryString did not did not return a valid queryString")
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
