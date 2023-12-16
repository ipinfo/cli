package lib

import (
	"fmt"
	"math"
	"math/big"
	"net"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/ipUtils"
	"github.com/spf13/pflag"
)

// CmdCalcFlags are flags expected by CmdCalc
type CmdCalcFlags struct {
	Help    bool
	NoColor bool
}

// Init initializes the common flags available to CmdCalc with sensible
func (f *CmdCalcFlags) Init() {
	_h := "see description in --help"
	pflag.BoolVarP(
		&f.Help,
		"help", "h", false,
		"show help.",
	)
	pflag.BoolVar(
		&f.NoColor,
		"nocolor", false,
		_h,
	)
}

// prec Function to return precedence of operators
func prec(s string) int {
	if s == "^" {
		return 3
	} else if (s == "/") || (s == "*") {
		return 2
	} else if (s == "+") || (s == "-") {
		return 1
	} else {
		return -1
	}
}

// isFloat Function to check if string is a float
func isFloat(str string) bool {
	pattern := `^[-+]?\d+(\.\d+)?$`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

// InfixToPostfix Function to convert infix expression to postfix expression using a stack based algorithm
func InfixToPostfix(infix []string) []string {
	var postfixStack Stack
	var postfix []string

	for _, token := range infix {
		if isOperator(token) {
			for {
				topOfStack, isEmpty := postfixStack.Pop()
				if isEmpty || prec(token) > prec(topOfStack) {
					postfixStack.Push(topOfStack)
					break
				}
				postfix = append(postfix, topOfStack)
			}
			postfixStack.Push(token)
		} else if token == "(" {
			postfixStack.Push(token)
		} else if token == ")" {
			for {
				topOfStack, _ := postfixStack.Pop()
				if topOfStack == "(" {
					break
				}
				postfix = append(postfix, topOfStack)
			}
		} else {
			postfix = append(postfix, token)
		}
	}

	// Pop all the remaining elements from the stack
	topOfStack, isEmpty := postfixStack.Pop()
	for !isEmpty {
		postfix = append(postfix, topOfStack)
		topOfStack, isEmpty = postfixStack.Pop()
	}
	return postfix
}

// EvaluatePostfix Function to evaluate postfix expression using a stack based algorithm
func EvaluatePostfix(postfix []string) (*big.Float, error) {
	// Precision for parsing string to big.Float
	var precision uint = 10000
	var postfixStack Stack
	for _, el := range postfix {
		// if operand, push it onto the stack.
		if el == "" {
			continue
		}
		if isFloat(el) || ipUtils.StrIsIPv4Str(el) || ipUtils.StrIsIPv6Str(el) {
			postfixStack.Push(el)
			continue
		}

		// if operator pop two elements off of the stack.
		strNum1, isEmpty := postfixStack.Pop()
		if isEmpty {
			return big.NewFloat(0), ipUtils.ErrInvalidInput
		}
		num1, _, _ := big.ParseFloat(strNum1, 10, precision, big.ToZero)

		strNum2, isEmpty := postfixStack.Pop()
		if isEmpty {
			return big.NewFloat(0), ipUtils.ErrInvalidInput
		}
		num2, _, _ := big.ParseFloat(strNum2, 10, precision, big.ToZero)

		operator := el
		result := new(big.Float)

		switch {
		case operator == "+":
			result = result.Add(num2, num1)
		case operator == "-":
			result = result.Sub(num2, num1)
		case operator == "*":
			result = result.Mul(num2, num1)
		case operator == "/":
			// Check for division by zero
			if num1.Cmp(big.NewFloat(0)) == 0 {
				return big.NewFloat(0), ipUtils.ErrInvalidInput
			}
			result = new(big.Float).Quo(num2, num1)
		case operator == "^":
			// Using Float64() to convert big.Float to float64
			// because big.Float does not have a equivalent function
			// for math.Pow() which accepts big.Float
			num1F64, _ := num1.Float64()
			num2F64, _ := num2.Float64()
			res := math.Pow(num2F64, num1F64)
			result = new(big.Float).SetPrec(precision).SetFloat64(res)
		default:
			return big.NewFloat(0), ipUtils.ErrInvalidInput
		}
		strResult := result.Text('f', 50)
		postfixStack.Push(strResult)
	}

	strTop, _ := postfixStack.Pop()
	top, _, _ := big.ParseFloat(strTop, 10, precision, big.ToZero)

	return top, nil
}

// isOperator Function to check if token is an operator
func isOperator(token string) bool {
	operators := map[string]bool{"+": true, "-": true, "*": true, "/": true, "^": true}
	_, isOperator := operators[token]
	return isOperator
}

// translateToken Function to translate token to decimal i.e. convert ipv4, ipv6 to decimal
func translateToken(tempToken string, tokens []string) ([]string, error) {
	if tempToken == "" {
		return tokens, nil
	}

	if isFloat(tempToken) {
		tokens = append(tokens, tempToken)
	} else if ipUtils.StrIsIPv4Str(tempToken) {
		// Convert ipv4 to decimal then append to tokens
		ip := net.ParseIP(tempToken)
		decimalIP := ipUtils.IP4toInt(ip)
		res := strconv.FormatInt(decimalIP, 10)
		tokens = append(tokens, res)

	} else if ipUtils.StrIsIPv6Str(tempToken) {
		ip := net.ParseIP(tempToken)
		decimalIP := ipUtils.IP6toInt(ip)
		tokens = append(tokens, decimalIP.String())
	} else {
		return []string{}, ipUtils.ErrInvalidInput
	}
	return tokens, nil
}

func isValidPartOfOperand(char rune) bool {
	validChars := `^[0-9a-fA-F:\.]*$`
	validCharsRegx := regexp.MustCompile(validChars)
	return validCharsRegx.MatchString(string(char))
}

// TokenizeInfix Function to tokenize infix expression
func TokenizeInfix(infix string) ([]string, error) {
	var tokens []string
	var err error
	var prevCharIsPartOfOperand bool
	var prevCharIsClosingBracket bool

	infix = "(" + infix + ")"
	tempToken := ""
	for _, char := range infix {
		opchar := string(char)
		if isValidPartOfOperand(char) {
			tempToken = tempToken + opchar
			prevCharIsPartOfOperand = true
			// If previous character was closing bracket and current character is part of operand
			// then insert a '*' between them
			// For example: (3+4)2 should be (3+4)*2
			if prevCharIsClosingBracket {
				tokens = append(tokens, "*")
			}
			prevCharIsClosingBracket = false

		} else if char == '(' || char == ')' || isOperator(opchar) {
			tokens, err = translateToken(tempToken, tokens)
			if err != nil {
				return []string{}, err
			}
			// If previous character was part of operand and current character is '('
			// then insert a '*' between them
			// For example: 2(3+4) should be 2*(3+4)
			if prevCharIsPartOfOperand && char == '(' {
				tokens = append(tokens, "*")
			}

			tokens = append(tokens, opchar)
			tempToken = ""

			prevCharIsPartOfOperand = false
			prevCharIsClosingBracket = char == ')'
		}
	}
	tokens = append(tokens, tempToken)
	return tokens, nil
}

// IsInvalidInfix Function to check if infix expression is valid
func IsInvalidInfix(expression string) bool {
	validChars := `^[0-9a-fA-F:\.\+\-\*\^\(\)\/ ]*$`
	validCharsRegx := regexp.MustCompile(validChars)

	var PrevChar rune
	var colonCount int
	for _, char := range expression {
		if isOperator(string(char)) && isOperator(string(PrevChar)) ||
			char == ')' && isOperator(string(PrevChar)) {
			return true
		}
		if char == '.' && PrevChar == '.' {
			return true
		}
		if char == ':' {
			colonCount++
			if colonCount > 2 {
				return true
			}
		} else {
			colonCount = 0
		}
		PrevChar = char
	}

	if isOperator(string(PrevChar)) || PrevChar == '.' {
		return true
	}

	return !validCharsRegx.MatchString(expression) || !isBalanced(expression)
}

// isBalanced Function to check if parentheses are balanced
func isBalanced(input string) bool {
	var postfixStack Stack
	for _, char := range input {
		if char == '(' {
			postfixStack.Push("(")
		} else if char == ')' {
			if postfixStack.IsEmpty() {
				return false
			}
			postfixStack.Pop()
		}
	}
	return postfixStack.IsEmpty()
}

// digitsAfterDecimal Function to count the number of non-zero digits after the decimal point
func digitsAfterDecimal(float big.Float) int {
	// Initially allowing 100 digits after decimal
	str := float.Text('f', 100)
	decimalIndex := strings.Index(str, ".")

	// Start counting the digits after the decimal point.
	count := 0
	for i := len(str) - 1; i > decimalIndex; i-- {
		if str[i] == '0' {
			count++
		} else {
			break
		}
	}

	return len(str) - (decimalIndex + 1) - count
}

// CmdCalc Function is the handler for the "calc" command.
func CmdCalc(f CmdCalcFlags, args []string, printHelp func()) error {
	if len(args) == 0 || f.Help {
		printHelp()
		return nil
	}

	if f.NoColor {
		color.NoColor = true
	}

	infix := args[0]
	if IsInvalidInfix(infix) {
		return ipUtils.ErrInvalidInput
	}

	tokens, err := TokenizeInfix(infix)
	if err != nil {
		return err
	}

	postfix := InfixToPostfix(tokens)
	result, err := EvaluatePostfix(postfix)
	if err != nil {
		return err
	}

	precision := digitsAfterDecimal(*result)
	resultStr := result.Text('f', precision)
	fmt.Println(resultStr)
	return nil
}
