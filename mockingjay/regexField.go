package mockingjay

import (
	"fmt"
	"regexp"
	"log"
)

// RegexYAML allows you to work with regex fields in YAML
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

// UnmarshalYAML will unhmarshal a YAML field into regexp
func (r *RegexField) UnmarshalJSON(data []byte) error {
	str := string(data)[1:len(data)-1]
	log.Println("Compiling", str)

	reg, err := regexp.Compile(str)

	if err != nil {
		return err
	}
	r.Regexp = reg
	return nil
}

// MarshalJSON returns a string for the regex
func (r *RegexField) MarshalJSON() ([]byte, error) {
	asString := fmt.Sprintf(`"%s"`, r.Regexp.String())
	return []byte(asString), nil
}
