package mockingjay

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
)

type MJDecoder interface {
	Decode(interface{}) error
}

type MJDecoderFunc func(interface{}) error

func (m MJDecoderFunc) Decode(target interface{}) error {
	return m(target)
}

func mjYAMLDecoder(data io.ReadCloser) MJDecoderFunc {
	return func(out interface{}) error {
		defer data.Close()
		stuff, err := ioutil.ReadAll(data)

		if err != nil {
			return fmt.Errorf("problem reading yaml data from source %v", err)
		}

		return yaml.Unmarshal(stuff, out)
	}
}
