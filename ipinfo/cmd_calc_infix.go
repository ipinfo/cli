package main

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"net"
	"os"
	"regexp"
	"strconv"
)

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

// Function to return precedence of operators
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

func isFloat(str string) bool {
	pattern := `^[-+]?\d+(\.\d+)?$`

	// Compile the regular expression pattern.
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(str)
}

func infixToPostfix(infix []string) []string {
	var sta Stack
	var postfix []string

	for _, token := range infix {
		if isOperator(token) {
			for !sta.IsEmpty() && prec(token) <= prec(sta.Top()) {
				postfix = append(postfix, sta.Top())
				sta.Pop()
			}
			sta.Push(token)
		} else if token == "(" {
			sta.Push(token)
		} else if token == ")" {
			for sta.Top() != "(" {
				postfix = append(postfix, sta.Top())
				sta.Pop()
			}
			sta.Pop()
		} else {
			postfix = append(postfix, token)
		}
	}
	// Pop all the remaining elements from the stack
	for !sta.IsEmpty() {
		postfix = append(postfix, sta.Top())
		sta.Pop()
	}
	return postfix
}

func evaluatePostfix(postfix []string) (*big.Float, error) {
	var sta Stack
	for _, el := range postfix {
		// if operand, push it onto the stack.
		if el == "" {
			continue
		}
		if isFloat(el) || isIPv4Address(el) || isIPv6Address(el) {
			sta.Push(el)
			continue
		}

		// if operator pop two elements off of the stack.
		var num1 big.Float
		strNum1 := sta.Top()
		_, success := num1.SetString(strNum1)

		if !success {
			fmt.Println("Error: Failed to convert the num1 to big.Int")
			return big.NewFloat(0), nil
		}
		sta.Pop()

		var num2 big.Float
		strNum2 := sta.Top()
		_, success = num2.SetString(strNum2)

		if !success {
			fmt.Println("Error: Failed to convert the num2 to big.Int:", strNum2)
			return big.NewFloat(0), nil
		}
		sta.Pop()
		operator := el

		result := new(big.Float)

		switch {
		case operator == "+":
			//fmt.Println("Adding")
			result = result.Add(&num2, &num1)

		case operator == "-":
			//fmt.Println("Subtracting")
			result = result.Sub(&num2, &num1)

		case operator == "*":
			//fmt.Println("Multiplying")
			result = result.Mul(&num2, &num1)
		case operator == "/":
			//fmt.Println("Dividing")
			result = new(big.Float).Quo(&num2, &num1)

		case operator == "^":
			result1, _ := num1.Float64()
			result2, _ := num2.Float64()

			res := math.Pow(result2, result1)
			result = new(big.Float).SetPrec(64).SetFloat64(res)

		default:
			fmt.Println("invalid operator: ", operator)
		}

		strResult := result.String()
		sta.Push(strResult)
	}

	strTop := sta.Top()
	sta.Pop()

	var top = new(big.Float)
	_, success := top.SetString(strTop)

	if !success {
		fmt.Println("Error: Failed to convert the string to big.Int")
		return big.NewFloat(0), nil
	}
	return top, nil
}

func isOperator(token string) bool {
	operators := map[string]bool{"+": true, "-": true, "*": true, "/": true, "^": true /* add other operators here */}
	_, isOperator := operators[token]
	return isOperator
}

func translateToken(tempToken string, tokens []string) ([]string, error) {
	var err error = nil

	if tempToken == "" {
		return tokens, nil
	}

	if isFloat(tempToken) {
		tokens = append(tokens, tempToken)
	} else if isIPv4Address(tempToken) {
		// convert ipv4 to decimal then append to tokens
		ip := net.ParseIP(tempToken)
		if ip == nil {
			err = errors.New("invalid IPv4 address: '" + tempToken + "'")
		}
		decimalIP := IP4toInt(ip)
		res := strconv.FormatInt(decimalIP, 10)
		tokens = append(tokens, res)

	} else if isIPv6Address(tempToken) {
		ip := net.ParseIP(tempToken)
		if ip == nil {
			fmt.Println("Invalid IPv6 address")
			err = errors.New("invalid IPv6 address: '" + tempToken + "'")
		}
		decimalIP := IP6toInt(ip)
		tokens = append(tokens, decimalIP.String())
	} else {
		err = errors.New("invalid expression")
	}
	return tokens, err
}

func tokeinzeExp(expression string) ([]string, error) {
	var tokens []string
	var err error

	expression = "(" + expression + ")"
	tempToken := ""
	for _, char := range expression {
		opchar := string(char)
		if isFloat(opchar) || opchar == "." || opchar == ":" {
			tempToken = tempToken + opchar
		} else if char == '(' || char == ')' || isOperator(opchar) {
			tokens, err = translateToken(tempToken, tokens)
			if err != nil {
				return []string{}, err
			}
			tokens = append(tokens, opchar)
			tempToken = ""
		}
	}
	tokens = append(tokens, tempToken)
	return tokens, nil
}

func IP6toInt(IPv6Address net.IP) *big.Int {
	IPv6Int := big.NewInt(0)
	IPv6Int.SetBytes(IPv6Address)
	return IPv6Int
}
func IP4toInt(IPv4Address net.IP) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(IPv4Address.To4())
	return IPv4Int.Int64()
}

func isIPv4Address(expression string) bool {
	// Define the regular expression pattern for matching IPv4 addresses
	ipV4Pattern := `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`

	// Compile the regular expression
	ipV4Regex := regexp.MustCompile(ipV4Pattern)

	// Use the MatchString function to check if the expression matches the IPv4 pattern
	return ipV4Regex.MatchString(expression)
}

func isIPv6Address(expression string) bool {
	// Define the regular expression pattern for matching IPv6 addresses
	ipV6Pattern := `^(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))$`

	// Compile the regular expression
	ipV6Regex := regexp.MustCompile(ipV6Pattern)

	// Use the MatchString function to check if the expression matches the IPv6 pattern
	return ipV6Regex.MatchString(expression)
}

func isInvalid(expression string) bool {
	validChars := `^[0-9:\.\+\-\*\^\(\)\/ ]*$`
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

// Function to check if parentheses are balanced
func isBalanced(input string) bool {
	var sta Stack
	for _, char := range input {
		if char == '(' {
			sta.Push("(")
		} else if char == ')' {
			if sta.IsEmpty() {
				return false
			}
			sta.Pop()
		}
	}
	return sta.IsEmpty()
}

func cmdCalcInfix() (string, error) {
	// infix := "2+3*(2^3-5)^(2+1*2)-4"
	cmd := ""
	if len(os.Args) > 2 {
		cmd = os.Args[2]
	}

	if isInvalid(cmd) {
		return "", errors.New("invalid expression")
	}

	tokens, err := tokeinzeExp(cmd)

	if err != nil {
		return "", err
	}

	postfix := infixToPostfix(tokens)

	result, err := evaluatePostfix(postfix)

	if err != nil {
		return "", err
	}

	return result.Text('f', 0), nil
}
