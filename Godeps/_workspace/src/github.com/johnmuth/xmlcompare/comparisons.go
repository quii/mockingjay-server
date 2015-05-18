// Package xmlcompare provides functions for comparing bits of XML.
package xmlcompare

import (
	"fmt"
	"reflect"

	"github.com/clbanning/mxj"
)

// IsCompatible compares two XML strings and returns true if the second has all the same element names and value types as the first.
// The elements don't have to be in the same order.
// The second XML string can have additional elements not present in the first.
func IsCompatible(a, b string) (compatible bool, err error) {
	aMap, err := mxj.NewMapXml([]byte(a), true)
	if err != nil {
		return
	}
	bMap, err := mxj.NewMapXml([]byte(b), true)
	if err != nil {
		return
	}
	return isStructurallyTheSame(aMap, bMap)
}

func isStructurallyTheSame(a, b map[string]interface{}) (compatible bool, err error) {
	for keyInA, v := range a {
		switch v.(type) {
		case string:
			_, bIsString := b[keyInA].(string)
			if bIsString {
				compatible = true
			}
		case float64:
			_, bIsFloat := b[keyInA].(float64)
			if bIsFloat {
				compatible = true
			}
		case map[string]interface{}:
			bMap, bIsMap := b[keyInA].(map[string]interface{})
			if bIsMap {
				for vKey, vValue := range v.(map[string]interface{}) {
					if elementValueType(vValue) != elementValueType(bMap[vKey]) {
						return
					}
				}
				compatible = true
			}
		default:
			err = fmt.Errorf("Unmatched datatype in XML found, got a %v", reflect.TypeOf(v))
		}
	}
	return
}

func elementValueType(elem interface{}) (valueType reflect.Type) {
	switch elem.(type) {
	case map[string]interface{}:
		valueType = reflect.TypeOf(elem.(map[string]interface{})["#text"])
	default:
		valueType = reflect.TypeOf(elem)
	}
	return
}
