package mockingjay

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"
	"encoding/json"
)

type testRegexDataType struct {
	Regex *RegexField
}

func TestItCanUnmarshalRegexFromYAML(t *testing.T) {
	rawYAML := `regex: \/hello\/[a-z]+`

	var d testRegexDataType
	err := yaml.Unmarshal([]byte(rawYAML), &d)

	assert.Nil(t, err)
	assert.NotNil(t, d.Regex)
}

func TestItReturnsErrorFoorInvalidRegex(t *testing.T) {
	rawYAML := `regex: "//!"!"£11\\/la%%\\/[a-z]+"`

	var d testRegexDataType
	err := yaml.Unmarshal([]byte(rawYAML), &d)

	assert.Error(t, err)
}

func TestItCanUnmarshalRegexFromJSON(t *testing.T) {
	rawJSON := `{"regex": "\/hello\/[a-z]+"}`

	var d testRegexDataType
	err := json.Unmarshal([]byte(rawJSON), &d)

	assert.Nil(t, err)
	assert.NotNil(t, d.Regex)
}
