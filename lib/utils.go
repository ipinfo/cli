package lib

// validateWithFunctions validates a slice of strings using a set of custom validation functions.
// It checks whether all elements in the input slice satisfy all the given validation functions.
// If any element fails validation for any of the functions, it returns false.
// Otherwise, it returns true.
//
//	Parameters:
//	s: A slice of strings to be validated.
//	fns: A slice of functions that take a string argument and return a boolean value. Each function represents a validation rule.
//
//	Returns:
//	true: If all elements in the input slice satisfy all the given validation functions.
//	false: If any element in the input slice fails validation for any of the functions.
func validateWithFunctions(s []string, fns []func(string) bool) bool {
	for _, fn := range fns {
		for _, str := range s {
			if !fn(str) {
				return false
			}
		}
	}
	return true
}

func StrContainsMultipleIPs(ipsStr []string) bool {
	numberOfIPs := 0
	for _, ipStr := range ipsStr {
		if StrIsIPStr(ipStr) {
			numberOfIPs++
		}
		if numberOfIPs > 1 {
			return true
		}
	}
	return false
}

func StrContainsMultipleASNs(asnsStr []string) bool {
	numberOfASNs := 0
	for _, asnStr := range asnsStr {
		if StrIsASNStr(asnStr) {
			numberOfASNs++
		}
		if numberOfASNs > 1 {
			return true
		}
	}
	return false
}
