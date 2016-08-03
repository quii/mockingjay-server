// Package jsonequaliser tries to determine if two JSON strings (A and B) are compatible with each other in terms of having the same fields matched to the same types. It tries to reduce false negatives by ignoring extra fields in B and will only require evidence of compatability. For example if A has an array with 3 objects, B just needs one.
package jsonequaliser
