package jsonequaliser

import "encoding/json"

func getJSONNodeFromString(data string) (node jsonNode, err error) {
	node = make(map[string]interface{})
	if err = json.Unmarshal([]byte(data), &node); err != nil {

		// Could be a top level array, in which case lets take the first item from it
		var anArr []jsonNode
		if err = json.Unmarshal([]byte(data), &anArr); err != nil {
			return
		}

		if len(anArr) > 0 {
			node = anArr[0]
		} else {
			err = new(emptyJSONArrayError)
		}

	}
	return
}

type emptyJSONArrayError struct{}

func (e *emptyJSONArrayError) Error() string {
	return "Empty arrays are not suitable for comparison"
}
