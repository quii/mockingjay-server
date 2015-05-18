package xmlcompare

import (
	"testing"
)

const baseXML = `<human><firstname>chris</firstname><lastname>james</lastname><age>30</age></human>`
const extraAttributeXML = `<human><firstname really="yes">chris</firstname><lastname>james</lastname><age>30</age></human>`
const compatibleXML = `<human><firstname>christopher</firstname><lastname>james</lastname><age>15</age></human>`
const differentOrderXML = `<human><age>30</age><firstname>chris</firstname><lastname>james</lastname></human>`
const differentElementNamesXML = `<wildebeest><name>Blue</name></wildebeest>`
const differentValueTypeXML = `<human><firstname>kristofferson</firstname><lastname>james</lastname><age>old</age></human>`
const invalidXML = `not even xml`

func TestIdenticalXMLIsCompatible(t *testing.T) {
	assertCompatible(t, baseXML, baseXML)
}

func TestSimpleXML(t *testing.T) {
	assertCompatible(t, `<foo>bar</foo>`, `<foo>baz</foo>`)
}

func TestSimpleXMLNumeric(t *testing.T) {
	assertCompatible(t, `<foo>1</foo>`, `<foo>2</foo>`)
}

func TestXMLWithSameElementNamesAndValueTypesIsCompatible(t *testing.T) {
	assertCompatible(t, baseXML, compatibleXML)
}

func TestXMLWithSameElementsInDifferentOrderIsCompatible(t *testing.T) {
	assertCompatible(t, baseXML, differentOrderXML)
}

func TestXMLWithAddedAttributesIsCompatible(t *testing.T) {
	assertCompatible(t, baseXML, extraAttributeXML)
}

func TestXMLWithDifferentElementNamesIsIncompatible(t *testing.T) {
	assertIncompatible(t, baseXML, differentElementNamesXML)
}

func TestXMLWithDifferentValueTypesIsIncompatible(t *testing.T) {
	assertIncompatible(t, baseXML, differentValueTypeXML)
}

func TestFirstInvalidXMLIsIncompatibleAndReturnsError(t *testing.T) {
	assertIncompatibleAndError(t, invalidXML, baseXML)
}

func TestSecondInvalidXMLIsIncompatibleAndReturnsError(t *testing.T) {
	assertIncompatibleAndError(t, baseXML, invalidXML)
}

func assertCompatible(t *testing.T, a, b string) {
	if compatible, err := IsCompatible(a, b); !compatible || err != nil {
		t.Errorf("%s should be compatible with %s (err = %v)", a, b, err)
	}
}

func assertIncompatible(t *testing.T, a, b string) {
	if compatible, err := IsCompatible(a, b); compatible || err != nil {
		t.Errorf("%s should not be compatible with %s (err = %v)", a, b, err)
	}
}

func assertIncompatibleAndError(t *testing.T, a, b string) {
	if compatible, err := IsCompatible(a, b); compatible || err == nil {
		t.Errorf("%s should not be compatible with %s and it should return an error", a, b)
	}
}
