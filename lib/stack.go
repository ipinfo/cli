package lib

// Stack represents a stack data structure that holds a collection of strings.
// The elements in the stack are managed using the Last-In-First-Out (LIFO) principle,
// which means that the last element added to the stack will be the first one to be removed.
type Stack []string

// IsEmpty checks if the stack is empty.
// It returns true if the stack contains no elements (i.e., it is empty), and false otherwise.
func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

// Push adds a new value onto the stack.
// The new value (given as a string) is appended to the top of the stack, effectively making it the new top element.
func (st *Stack) Push(str string) {
	*st = append(*st, str) //Simply append the new value to the end of the stack
}

// Pop removes the top element from the stack and returns it. If the stack is empty,
// it returns an empty string and true indicating that the stack is empty.
// Otherwise, it returns the top element of the stack and false.
// The stack is modified in-place by removing the top element.
func (st *Stack) Pop() (string, bool) {
	if st.IsEmpty() {
		return "", true
	} else {
		index := len(*st) - 1   // Get the index of top most element.
		element := (*st)[index] // Index onto the slice and obtain the element.
		*st = (*st)[:index]     // Remove it from the stack by slicing it off.
		return element, false
	}
}
