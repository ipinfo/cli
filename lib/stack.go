package lib

// Stack type
type Stack []string

// IsEmpty check if stack is empty
func (st *Stack) IsEmpty() bool {
	return len(*st) == 0
}

// Push a new value onto the stack
func (st *Stack) Push(str string) {
	*st = append(*st, str) //Simply append the new value to the end of the stack
}

// Pop Remove top element of stack. Return false if stack is empty.
func (st *Stack) Pop() bool {
	if st.IsEmpty() {
		return false
	} else {
		index := len(*st) - 1 // Get the index of top most element.
		*st = (*st)[:index]   // Remove it from the stack by slicing it off.
		return true
	}
}

// Top Return top element of stack. Return false if stack is empty.
func (st *Stack) Top() string {
	if st.IsEmpty() {
		return ""
	} else {
		index := len(*st) - 1   // Get the index of top most element.
		element := (*st)[index] // Index onto the slice and obtain the element.
		return element
	}
}
