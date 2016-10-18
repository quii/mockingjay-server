package mockingjay

import (
	"fmt"
	"regexp"
)

// RegexField allows you to work with regex fields in YAML
type RegexField struct {
	*regexp.Regexp
}

// UnmarshalYAML will unhmarshal a YAML field into regexp
func (r *RegexField) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var stringFromYAML string
	err := unmarshal(&stringFromYAML)
	if err != nil {
		return err
	}
	reg, err := regexp.Compile(stringFromYAML)
	if err != nil {
		return err
	}
	r.Regexp = reg
	return nil
}

// UnmarshalJSON will unhmarshal a JSON field into regexp
func (r *RegexField) UnmarshalJSON(data []byte) error {
	str := string(data)[1 : len(data)-1]
	reg, err := regexp.Compile(str)

	if err != nil {
		return err
	}
	r.Regexp = reg
	return nil
}

// MarshalJSON returns a string for the regex
func (r *RegexField) MarshalJSON() ([]byte, error) {
	asString := fmt.Sprintf(`"%v"`, r.Regexp.String())
	return []byte(asString), nil
}

// MarshalYAML returns the string of the regex
func (r *RegexField) MarshalYAML() (interface{}, error) {
	return r.Regexp.String(), nil
}
