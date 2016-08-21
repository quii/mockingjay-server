package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItValidatesRequests(t *testing.T) {

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
