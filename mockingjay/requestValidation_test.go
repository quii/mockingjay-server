package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItValidatesRequests(t *testing.T) {

	spaceyHeaders := make(map[string]string)
	spaceyHeaders["hello world"] = "l o l"

	aForm := make(map[string]string)
	aForm["name"] = "Hudson"

	tests := []struct {
		Description   string
		Request       Request
		ExpectedError error
	}{
		{
			Description: "Empty URIs are not valid",
			Request: Request{
				URI:    "",
				Method: "POST",
			},
			ExpectedError: errEmptyURI,
		},
		{
			Description: "Empty methods are not valid",
			Request: Request{
				URI:    "/",
				Method: "",
			},
			ExpectedError: errEmptyMethod,
		},
		{
			Description: "Headers cant have spaces",
			Request: Request{
				URI:     "/",
				Method:  "POST",
				Headers: spaceyHeaders,
			},
			ExpectedError: errBadHeaders,
		},
		{
			Description: "Cant have both form and a body",
			Request: Request{
				URI:    "/",
				Method: "POST",
				Body:   "hi",
				Form:   aForm,
			},
			ExpectedError: errBodyAndForm,
		},
	}
	for _, test := range tests {
		assert.Equal(t, test.ExpectedError, test.Request.errors(), test.Description)
	}

	validRequest := Request{
		URI:    "/",
		Method: "POST",
	}

	assert.Nil(t, validRequest.errors())

}
