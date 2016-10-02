package mockingjay

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"encoding/json"
	"gopkg.in/yaml.v2"
	"regexp"
	"strings"
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

func TestItCanMarshalBackToJSON(t *testing.T) {
	regexURI, err := regexp.Compile(`\/hello\/[a-z]+`)

	if err != nil {
		t.Fatal("Cant compile regex", err)
	}

	field := testRegexDataType{&RegexField{regexURI}}

	data, err := json.Marshal(field)

	if err != nil {
		t.Error("Couldn't marshal into JSON", err)
	}

	assert.Equal(t, `{"Regex":"\/hello\/[a-z]+"}`, string(data))
}

func TestItCanMarshalBackToYAML(t *testing.T) {
	regexURI, err := regexp.Compile(`\/hello\/[a-z]+`)

	if err != nil {
		t.Fatal("Cant compile regex", err)
	}

	field := testRegexDataType{&RegexField{regexURI}}

	data, err := yaml.Marshal(field)

	if err != nil {
		t.Error("Couldn't marshal into JSON", err)
	}

	assert.Equal(t, `regex: \/hello\/[a-z]+`, strings.Replace(string(data), "\n", "", -1))
}
