package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItValidatesResponses(t *testing.T) {
	tests := []struct {
		Description string
		Response    response
	}{
		{
			Description: "Response codes must be > 99",
			Response: response{
				Code: 99,
			},
		},
		{
			Description: "Response codes must be < 600 ish",
			Response: response{
				Code: 600,
			},
		},
	}

	for _, test := range tests {
		assert.False(t, test.Response.isValid(), test.Description)
	}

}
