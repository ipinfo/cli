package lib

// ReverseSliceString reverses the order of elements inside of a string slice.
func ReverseSliceString(s []string) {
	last := len(s) - 1
	for i := 0; i < len(s)/2; i++ {
		s[i], s[last-i] = s[last-i], s[i]
	}
}
