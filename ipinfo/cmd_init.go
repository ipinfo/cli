package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
	"golang.org/x/term"
)

var completionsInit = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-t":         predict.Nothing,
		"--token":    predict.Nothing,
		"--no-check": predict.Nothing,
		"-h":         predict.Nothing,
		"--help":     predict.Nothing,
	},
}

func printHelpInit() {
	fmt.Printf(
		`Usage: %s init [<opts>] [<token>]
Examples:
	# Login command with token flag.
	$ init --token <token>

    # Authentication without token flag.
	$ init

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

type meResponse struct {
	Error string `json:"error"`
}

func cmdInit() error {
	var opt int
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
		printHelpInit()
		return nil
	}

	// allow only flag or arg for token but not both.
	if fTok != "" && len(args) > 0 {
		return errors.New("ambiguous token input source")
	}

	// get token, from flag or command line.
	// if it exists, we'll exit early as it's an implicit login.
	tok := fTok
	if len(args) > 0 {
		tok = args[0]
	}
	if tok != "" {
		if !fNoCheck {
			if err := checkValidity(tok); err != nil {
				return fmt.Errorf("could not confirm if token is valid: %w", err)
			}
		}

		// save token to file.
		gConfig.Token = tok
		if err := SaveConfig(gConfig); err != nil {
			return err
		}

		fmt.Println("done")

		return nil
	}

	fmt.Println("1) Enter an existing API token")
	fmt.Println("2) Create a new account")
	_, err := fmt.Scanf("%d", &opt)
	if err != nil {
		return fmt.Errorf("error reading input: %v", err)
	}
	if opt == 1 {
		newtoken, err := enterToken(tok)
		if err != nil {
			return fmt.Errorf(err.Error())
		}

		// check token validity.
		if !fNoCheck {
			if err := checkValidity(newtoken); err != nil {
				return fmt.Errorf("could not confirm if token is valid: %w", err)
			}
		}

		// save token to file.
		gConfig.Token = tok
		if err := SaveConfig(gConfig); err != nil {
			return err
		}

		fmt.Println("done")

		return nil
	} else if opt == 2 {
		res, err := http.Get("https://ipinfo.io/signup/cli")
		if err != nil {
			return err
		}
		defer res.Body.Close()

		// parse response.
		rawBody, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		body := &signupCli{}
		err = json.Unmarshal(rawBody, body)
		if err != nil {
			return err
		}

		_ = openURL(body.SignupURL)
		fmt.Println("If the link does not open, please go to this link to get your access token:")
		fmt.Println("")
		fmt.Printf("%v\n", body.SignupURL)
		fmt.Println("")
		fmt.Println("Press [Enter] when done if not automatically detected.")

		// Retrieving CLI token from signup URL.
		parsedUrl, err := url.Parse(body.SignupURL)
		if err != nil {
			fmt.Println("Error parsing URL:", err)
			return err
		}
		cliToken := parsedUrl.Query().Get("cli_token")
		if cliToken == "" {
			fmt.Println("CLI token not found in URL")
			return err
		}

		// Check if signup flow is completed.
		maxAttempts := 200
		count := 0
		for {
			count++
			res, err := http.Get("https://ipinfo.io/signup/cli/check?cli_token=" + cliToken)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
			defer res.Body.Close()

			if res.StatusCode == http.StatusOK {
				rawBody, err := io.ReadAll(res.Body)
				if err != nil {
					return err
				}
				body := &tokenCli{}
				err = json.Unmarshal(rawBody, body)
				if err != nil {
					return err
				}

				// save token to file.
				gConfig.Token = tok
				if err := SaveConfig(gConfig); err != nil {
					return err
				}

				fmt.Println("Account created successfully.")
				break
			}

			if count == maxAttempts {
				if _, err := fmt.Scanln(); err != nil {
					return fmt.Errorf("%v", err)
				}
			}

			time.Sleep(time.Second)
		}
	} else {
		fmt.Println("Invalid input.")
		return err
	}

	return nil
}

func checkValidity(tok string) error {
	fmt.Println("checking token...")
	tokenOk, err := isTokenValid(tok)
	if err != nil {
		return fmt.Errorf("could not confirm if token is valid: %w", err)
	}
	if !tokenOk {
		return fmt.Errorf("invalid token")
	}

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

	return tok, nil
}

func isTokenValid(tok string) (bool, error) {
	// make API req for true token validity.
	res, err := http.Get("https://ipinfo.io/me?token=" + tok)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	// parse response.
	me := &meResponse{}
	if err := json.NewDecoder(res.Body).Decode(me); err != nil {
		return false, err
	}

	// If no errors then me.Error should be empty
	return me.Error == "", nil
}

func openURL(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
