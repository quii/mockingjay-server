package mockingjay

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"gopkg.in/yaml.v2"
)

type testRegexDataType struct {
	Regex *RegexYAML
}

func TestItCanUnmarshalRegex(t *testing.T) {
	rawYAML := `regex: "\\/hello\\/[a-z]+"`

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
