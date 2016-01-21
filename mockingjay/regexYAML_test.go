package mockingjay

import (
	"testing"

	"gopkg.in/yaml.v2"
)

type testRegexDataType struct {
	Regex *RegexYAML
}

func TestItCanUnmarshalRegex(t *testing.T) {
	rawYAML := `regex: "\\/hello\\/[a-z]+"`

	var d testRegexDataType
	err := yaml.Unmarshal([]byte(rawYAML), &d)

	if err != nil {
		t.Error("Couldnt unmarshal regex from YAML", err)
	}

	if d.Regex == nil {
		t.Error("Regex was not extracted from YAML")
	}
}
