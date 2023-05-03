package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/pflag"
	"golang.org/x/term"
)

func printHelpInit() {
	fmt.Printf(
		`Usage: %s init [<opts>] [<token>]

Options:
  --token <tok>, -t <tok>
    token to login with.
    (this is potentially unsafe; let the CLI prompt you instead).
  --no-check
    disable checking if the token is valid or not.
    default: false.
  --help, -h
    show help.
`, progBase)
}

type signupCli struct {
	SignupURL string `json:"signupURL"`
}

type tokenCli struct {
	Token string `json:"token"`
}

func cmdInit() error {
	var num int
	var fTok string
	var fNoCheck bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to save.")
	pflag.BoolVar(&fNoCheck, "no-check", false, "disable checking if token is valid.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpInit()
		return nil
	}

	// get args without subcommand.
	args := pflag.Args()[1:]

	// only token arg allowed.
	if len(args) > 1 {
		return errors.New("invalid arguments")
	}

	// allow only flag or arg for token but not both.
	if fTok != "" && len(args) > 0 {
		return errors.New("ambiguous token input source")
	}

	// get token, from flag or command line.
	tok := fTok
	if len(args) > 0 {
		tok = args[0]
	}
	if tok != "" && !fNoCheck {
		if err := checkValidity(tok); err != nil {
			return fmt.Errorf("could not confirm if token is valid: %w", err)
		} else {
			return nil
		}
	}
	fmt.Printf("1) Enter an existing API token\n")
	fmt.Printf("2) Create a new account\n")
	_, err := fmt.Scanf("%d", &num)
	if err != nil {
		fmt.Println("Error reading input:", err)
		return err
	}
	if num == 1 {
		newtoken, err := enterToken(tok);
		if err != nil {
			return fmt.Errorf(err.Error())
		}
		if err := checkValidity(newtoken); err != nil {
			return fmt.Errorf("could not confirm if token is valid: %w", err)
		}
	} else if num == 2 {
		res, err := http.Get("https://ipinfo.io/signup/cli")
		if err != nil {
			return err
		}
		defer res.Body.Close()
		// parse response.
		msg := &signupCli{}
		if err := json.NewDecoder(res.Body).Decode(msg); err != nil {
			return err
		}
		fmt.Printf("%v\n",msg.SignupURL)
		var input string
		fmt.Println("Press Enter to open link:")
		fmt.Scanf("%s", &input)
		cmd := exec.Command("xdg-open", msg.SignupURL)
		err = cmd.Run()
		if err != nil {
			fmt.Println("Error opening link:", err)
			return err
		}
		parsedUrl, err := url.Parse(msg.SignupURL)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return err
		}
	
		uid := parsedUrl.Query().Get("uid")
		if uid == "" {
			fmt.Println("UID not found in URL")
			return err
		}

		// Check if signup flow is completed.
		maxAttempts := 200
		count := 0
		fmt.Printf("Configuring account ...\n")
		for {
			count++
			res, err := http.Get("https://ipinfo.io/signup/cli/check?uid=" + uid)
			if err != nil {
				return fmt.Errorf("%v",err)
			}
			defer res.Body.Close()
	
			if res.StatusCode == http.StatusOK {
				body := &tokenCli{}
				if err := json.NewDecoder(res.Body).Decode(body); err != nil {
					return err
				}
				if tok, err = enterToken(body.Token); err != nil {
					return fmt.Errorf(err.Error())
				}
				if err := checkValidity(tok); err != nil {
					return fmt.Errorf("could not confirm if token is valid: %w", err)
				}
				break
			}
	
			if count == maxAttempts {
				fmt.Println("Reached max attempts. Press Enter to retry or Ctrl+C to exit.")
				if _, err := fmt.Scanln(); err != nil {
					return fmt.Errorf("%v",err)
				}
				// reset the count if user chooses to retry
				count = 0
			}

			time.Sleep(time.Second)
		}
		os.Exit(0)
	} else {
		fmt.Println("Invalid input.")
		return err
	}

	return nil
}

func checkValidity(tok string) error {
	fmt.Println("logging in...")
	tokenOk, err := isTokenValid(tok)
	if err != nil {
		return fmt.Errorf("could not confirm if token is valid: %w", err)
	}
	if !tokenOk {
		return fmt.Errorf("invalid token")
	}

	// save token to file.
	gConfig.Token = tok
	if err := SaveConfig(gConfig); err != nil {
		return err
	}

	fmt.Println("done")

	return nil
}

func enterToken(tok string) (string, error) {
	for tok == "" {
		fmt.Printf("Enter token: ")
		tokraw, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}

		tok = string(tokraw[:])

		// exit if we have a token now.
		if tok != "" {
			break
		}

		fmt.Println("please enter a token")
	}

	return tok , nil
}
