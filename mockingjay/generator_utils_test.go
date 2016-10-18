package mockingjay

import (
	"fmt"
	"net/url"
	"testing"
	"testing/quick"
)

func TestRandomURL(t *testing.T) {

	f := func(n uint8) bool {
		fmt.Println(n)
		path := randomURL(int(n) + 1)
		_, err := url.Parse(path)
		return err == nil
	}

	if err := quick.Check(f, nil); err != nil {
		t.Error("randomUrl did not return a valid URL")
	}

}
