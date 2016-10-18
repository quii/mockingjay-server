package mockingjay

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItAddsNewLinesToYAML(t *testing.T) {
	badYAML := `- name: Test endpoint
  cdcdisabled: false
  request:
    uri: /hello
    method: GET
  response:
    code: 200
    body: '{"message": "hello, world"}'
    headers:
      content-type: text/json
- name: Test endpoint 2
  cdcdisabled: false
  request:
    uri: /world
    method: DELETE
  response:
    code: 200
    body: hello, world`

	expectedYAML := `
- name: Test endpoint
  cdcdisabled: false
  request:
    uri: /hello
    method: GET
  response:
    code: 200
    body: '{"message": "hello, world"}'
    headers:
      content-type: text/json

- name: Test endpoint 2
  cdcdisabled: false
  request:
    uri: /world
    method: DELETE
  response:
    code: 200
    body: hello, world
`

	assert.Equal(t, expectedYAML, string(addNewLinesToConfig([]byte(badYAML))))

}
