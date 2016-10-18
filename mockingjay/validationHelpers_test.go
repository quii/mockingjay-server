package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItSeesEmptyHeaderValuesAsInvalid(t *testing.T) {
	emptyValHeaders := make(map[string]string)

	emptyValHeaders["foo"] = ""

	assert.False(t, httpHeadersValid(emptyValHeaders), "Empty value for a header is invalid")
}
