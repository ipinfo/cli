package lib

type ValidatorFunc func(string) bool

func ValidateWithAnyValidator(strings []string, validators []ValidatorFunc) bool {
	if len(strings) == 0 || len(validators) == 0 {
		return false // Return false if either the input strings or validators are empty.
	}

	for _, str := range strings {
		valid := false // Flag to keep track of successful validation for the current string.
		for _, validator := range validators {
			if validator(str) {
				valid = true // At least one validator returned true for the current string.
				break        // No need to check further validators for this string.
			}
		}
		if !valid {
			return false // Return false if no validator returned true for the current string.
		}
	}

	return true // All strings passed at least one validation function.
}

func StrArrIsCombinationOfIPsAndASNs(strArr []string) bool {
	return len(strArr) >= 2 &&
		ValidateWithAnyValidator(strArr, []ValidatorFunc{StrIsIPStr, StrIsASNStr})
}
