package jsonequaliser

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// IsCompatible checks that two json strings are structurally the same so that they are compatible. The first string should be your "correct" json, if there are extra fields in B then they will still be seen as compatible
func IsCompatible(a, b string) (compatible bool, err error) {

	aMap := make(map[string]interface{})
	if err = json.Unmarshal([]byte(a), &aMap); err != nil {
		return
	}

	bMap := make(map[string]interface{})
	if err = json.Unmarshal([]byte(b), &bMap); err != nil {
		return
	}

	return isStructurallyTheSame(aMap, bMap)
}

func isStructurallyTheSame(a, b map[string]interface{}) (compatible bool, err error) {
	for keyInA, v := range a {

		if b[keyInA] == nil {
			return
		}

		switch v.(type) {
		case string:
			if _, isString := b[keyInA].(string); !isString {
				return
			}
		case bool:
			if _, isBool := b[keyInA].(bool); !isBool {
				return
			}
		case float64:
			if _, isFloat := b[keyInA].(float64); !isFloat {
				return
			}

		case interface{}:
			aArr, _ := a[keyInA].([]interface{})

			bArr, bIsArray := b[keyInA].([]interface{})

			if !bIsArray || len(bArr) == 0 {
				return
			}

			aLeaf, aIsMap := aArr[0].(map[string]interface{})
			bLeaf, bIsMap := bArr[0].(map[string]interface{})

			if aIsMap && bIsMap {
				return isStructurallyTheSame(aLeaf, bLeaf)
			} else if aIsMap && !bIsMap {
				return
			} else if reflect.TypeOf(aArr[0]) != reflect.TypeOf(bArr[0]) {
				return
			}
		default:
			err = fmt.Errorf("Unmatched type of json found, got a %v", reflect.TypeOf(v))
			return
		}
	}
	compatible = true
	return
}
