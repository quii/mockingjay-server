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

func contains(char rune, s string) bool {
	for _, r := range s {
		if r == char {
			return true
		}
	}
	return false
}
