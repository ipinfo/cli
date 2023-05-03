package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

const dbDownloadURL = "https://ipinfo.io/data/free/"

var completionsDownload = &complete.Command{
	Flags: map[string]complete.Predictor{
		"--asn":         predict.Nothing,
		"--country":     predict.Nothing,
		"--country-asn": predict.Nothing,
		"-f":            predict.Nothing,
		"--format":      predict.Nothing,
		"-t":            predict.Nothing,
		"--token":       predict.Nothing,
		"-h":            predict.Nothing,
		"--help":        predict.Nothing,
	},
}

func printHelpDownload() {
	fmt.Printf(
		`Usage: %s download [<opts>]

Description:
    Download the free ipinfo databases.

Examples:
    # Download country database in csv format.
    $ %[1]s download --country --csv

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --asn
      download the free ipinfo asn database. 
    --country
      download the free ipinfo country database. 
    --country-asn
      download the free ipinfo country asn database. 
    --help, -h
      show help.

Outputs:
    --format , -f <mmdb | json | csv>
    output format of the database file.
      mmdb (default.) => downloads the mmdb format database.
      json            => downloads the json format database.
      csv             => downloads the csv  format database.
`, progBase)
}

func cmdDownload() error {
	var fTok string
	var fFmt string
	var fAsn bool
	var fCountry bool
	var fCountryAsn bool
	var fHelp bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVar(&fAsn, "asn", false, "free asn database.")
	pflag.BoolVar(&fCountry, "country", false, "free country database.")
	pflag.BoolVar(&fCountryAsn, "country-asn", false, "free country asn database.")
	pflag.StringVarP(&fFmt, "format", "f", "mmdb", "the output format to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if pflag.NFlag() == 0 {
		printHelpDownload()
		return nil
	}

	if fHelp {
		printHelpDownload()
		return nil
	}

	var token string
	if fTok == "" {
		token = gConfig.Token
	}

	// require token for download.
	if token == "" {
		return errors.New("downloading requires a token; login via `ipinfo login` or pass the `--token` argument")
	}

	// check download format
	var format string
	switch fFmt {
	case "mmdb":
		format = "mmdb"
	case "csv":
		format = "csv.gz"
	case "json":
		format = "json.gz"
	default:
		return errors.New("unknown download format")
	}

	if fAsn {
		err := downloadDb("asn", format, token)
		if err != nil {
			return err
		}
	}
	if fCountry {
		err := downloadDb("country", format, token)
		if err != nil {
			return err
		}
	}
	if fCountryAsn {
		err := downloadDb("country_asn", format, token)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadDb(name, format, token string) error {
	url := fmt.Sprintf("%s%s.%s?token=%s", dbDownloadURL, name, format, token)

	// make API req to download the file.
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// create file.
	fileName := fmt.Sprintf("%s.%s", name, format)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// save file.
	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Database %s.%s saved successfully.", name, format)
	return nil
}
