package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func IPInputAction(
	inputs []string,
	stdin bool,
	ip bool,
	iprange bool,
	cidr bool,
	file bool,
	actionStdin func(string),
	actionRange func(string),
	actionCidr func(string),
	actionFile func(string),
) error {
	// Handle stdin
	if stdin {
		stat, _ := os.Stdin.Stat()

		isPiped := (stat.Mode() & os.ModeNamedPipe) != 0
		isTyping := (stat.Mode()&os.ModeCharDevice) != 0 && len(inputs) == 0

		if isTyping {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}

		if isPiped || isTyping || stat.Size() > 0 {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				input := strings.TrimSpace(scanner.Text())
				if input == "" {
					continue
				}
				actionStdin(input)
			}
		}
	}

	// Parse inputs
	for _, input := range inputs {
		if iprange && actionRange != nil {
			actionRange(input)
		}

		if ip && StrIsIPStr(input) {
			actionStdin(input)
		}

		if cidr && StrIsCIDRStr(input) {
			actionCidr(input)
		}

		if file && FileExists(input) && actionFile != nil {
			actionFile(input)
		}
	}

	return nil
}
