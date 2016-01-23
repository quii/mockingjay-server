package jsonequaliser

import (
	"fmt"
	"reflect"
)

type jsonNode map[string]interface{}

// IsCompatible checks that two json strings are structurally the same so that they are compatible. The first string should be your "correct" json, if there are extra fields in B then they will still be seen as compatible
func IsCompatible(a, b string) (compatible bool, err error) {

	aMap, err := getJSONNodeFromString(a)
	bMap, err := getJSONNodeFromString(b)

	if err != nil {
		return
	}

	return isStructurallyTheSame(aMap, bMap)
}

func isStructurallyTheSame(a, b jsonNode) (compatible bool, err error) {
	for jsonFieldName, v := range a {

		if fieldMissingIn(b, jsonFieldName) {
			return false, nil
		}

		if a[jsonFieldName] == nil {
			return true, nil
		}

		switch v.(type) {
		case string:
			if !isString(b, jsonFieldName) {
				return
			}
		case bool:
			if !isBool(b, jsonFieldName) {
				return
			}
		case float64:
			if !isFloat(b, jsonFieldName) {
				return
			}

		case interface{}:

			aArr, aIsArray := a[jsonFieldName].([]interface{})

			bArr, bIsArray := b[jsonFieldName].([]interface{})

			if aIsArray && len(aArr) == 0 {
				return true, nil
			}

			if !bIsArray && aIsArray || aIsArray && len(bArr) == 0 {
				return
			}

			var aLeaf, bLeaf jsonNode
			var aIsMap, bIsMap bool

			if aIsArray && bIsArray {
				aLeaf, aIsMap = aArr[0].(map[string]interface{})
				bLeaf, bIsMap = bArr[0].(map[string]interface{})
			} else {
				aLeaf, aIsMap = a[jsonFieldName].(map[string]interface{})
				bLeaf, bIsMap = b[jsonFieldName].(map[string]interface{})
			}

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
