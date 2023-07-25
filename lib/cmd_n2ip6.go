package lib

func CmdN2IP6(expression string) (string, error) {
	if IsInvalidInfix(expression) {
		return "", ErrInvalidInput
	}

	// n2ip6 also accepts an expression which is why the following
	// Steps are being done
	// Convert to postfix
	// If it is a single number and not an expression
	// The tokenization and evaluation would have no effect on the number

	// Tokenize
	tokens, err := TokenizeInfix(expression)
	if err != nil {
		return "", err
	}

	postfix := InfixToPostfix(tokens)

	// Evaluate the postfix expression
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		return "", err
	}

	// Convert to IP
	// Precision should be 0 i.e. number of digits after decimal
	// as ip cannot be derived from a float
	res, err := DecimalStrToIP(result.Text('f', 0), true)
	if err != nil {
		return "", err
	}

	return res.String(), nil
}
