package jsonequaliser

func fieldMissingIn(node jsonNode, field string) bool {
	_, exists := node[field]
	return !exists
}

func isString(node jsonNode, key string) bool {
	_, isString := node[key].(string)
	return isString
}

func isBool(node jsonNode, key string) bool {
	_, isString := node[key].(bool)
	return isString
}

func isFloat(node jsonNode, key string) bool {
	_, isString := node[key].(float64)
	return isString
}
