package tokener

type Tokener struct {
	quotes  []byte
	escaped bool
	fixed   string
	space   bool
}

// Visit visit a byte and update the state of the quotes.
// It returns true if the byte was quotes or escape character.
func (t *Tokener) Visit(b byte) {
	// Check space.
	if b == ' ' {
		if !t.escaped && !t.Quoted() {
			t.space = true
		}
	} else {
		t.space = false
	}

	// Check escaping
	if b == '\\' {
		t.escaped = !t.escaped
	} else {
		defer func() { t.escaped = false }()
	}

	// Check quotes.
	if !t.escaped && (b == '"' || b == '\'') {
		if t.Quoted() && t.quotes[len(t.quotes)-1] == b {
			t.quotes = t.quotes[:len(t.quotes)-1]
		} else {
			t.quotes = append(t.quotes, b)
		}
	}

	// If not quoted, insert escape before inserting space.
	if t.LastSpace() {
		t.fixed += "\\"
	}
	t.fixed += string(b)
}

func (t *Tokener) Escaped() bool {
	return t.escaped
}

func (t *Tokener) Quoted() bool {
	return len(t.quotes) > 0
}

func (t *Tokener) Fixed() string {
	return t.fixed
}

func (t *Tokener) Closed() string {
	fixed := t.fixed
	for i := len(t.quotes) - 1; i >= 0; i-- {
		fixed += string(t.quotes[i])
	}
	return fixed
}

func (t *Tokener) LastSpace() bool {
	return t.space
}
