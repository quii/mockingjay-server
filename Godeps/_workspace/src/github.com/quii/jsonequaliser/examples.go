package jsonequaliser

import "fmt"

// ExampleIsCompatible shows how the JSON is compatible, even when the data is different. JSONEqualiser just checks the type of the fields
func ExampleIsCompatible() {
	A := `[{"firstname": "chris", "lastname": "james", "age": 30}]`
	B := `[{"firstname": "Bob", "lastname": "Smith", "age": 25, "favourite-colour": "blue"}]`

	fmt.Println(IsCompatible(A, B))
	// Output: map[] <nil>
}

// ExampleIsIncompatible shows what happens when there is a problem. When consuming this lib you should check to see if the map is empty.
func ExampleIsIncompatible() {
	A := `{"foo": "bar"}`
	B := `{"foo": true}`

	fmt.Println(IsCompatible(A, B))
	// map[foo:Field is not a string in other JSON] <nil>
}
