package main

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

const dbDownloadURL = "https://ipinfo.io/data/free/"

var completionsDownload = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-c":         predict.Nothing,
		"--compress": predict.Nothing,
		"-f":         predict.Nothing,
		"--format":   predict.Nothing,
		"-t":         predict.Nothing,
		"--token":    predict.Nothing,
		"-h":         predict.Nothing,
		"--help":     predict.Nothing,
	},
}

func printHelpDownload() {
	fmt.Printf(
		`Usage: %s download [<opts>] <database> [<output>]

Description:
    Download the free ipinfo databases.

Examples:
    # Download country database in csv format.
    $ %[1]s download country -f csv > country.csv
    $ %[1]s download country-asn country_asn.mmdb

Databases:
    asn            free ipinfo asn database.
    country        free ipinfo country database.
    country-asn    free ipinfo country-asn database.

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

Outputs:
    --compress, -c
	save the file in compressed format.
	default: false.
    --format, -f <mmdb | json | csv>
     output format of the database file.
     default: mmdb.
`, progBase)
}

func cmdDownload() error {
	var fTok string
	var fFmt string
	var fZip bool
	var fHelp bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.StringVarP(&fFmt, "format", "f", "mmdb", "the output format to use.")
	pflag.BoolVarP(&fZip, "compress", "c", false, "compressed output.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	args := pflag.Args()[1:]
	if fHelp || len(args) > 2 || len(args) < 1 {
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

	// get download format and extension.
	var format string
	var fileExtension string
	switch strings.ToLower(fFmt) {
	case "mmdb":
		format = "mmdb"
		fileExtension = "mmdb"
	case "csv":
		format = "csv.gz"
		fileExtension = "csv"
	case "json":
		format = "json.gz"
		fileExtension = "json"
	default:
		return errors.New("unknown download format")
	}

	if fZip {
		fileExtension = fmt.Sprintf("%s.%s", fileExtension, "gz")
	}

	// download the db.
	switch strings.ToLower(args[0]) {
	case "asn":
		err := downloadDb("asn", format, token, fileExtension, fZip)
		if err != nil {
			return err
		}
	case "country":
		err := downloadDb("country", format, token, fileExtension, fZip)
		if err != nil {
			return err
		}
	case "country-asn":
		err := downloadDb("country_asn", format, token, fileExtension, fZip)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("database '%v' is invalid", args[0])
	}

	return nil
}

func downloadDb(name, format, token, fileExtension string, zip bool) error {
	url := fmt.Sprintf("%s%s.%s?token=%s", dbDownloadURL, name, format, token)

	// get file name.
	var fileName string
	if len(pflag.Args()) > 2 {
		fileName = pflag.Args()[2]
	} else {
		fileName = fmt.Sprintf("%s.%s", name, fileExtension)
	}

	// make API req to download the file.
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// if output not terminal unzip and write to stdout.
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		err := unzipWrite(os.Stdout, res.Body)
		if err != nil {
			return err
		}
	} else {
		// create file.
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// save compressed file.
		if zip {
			if format == "mmdb" {
				writer := gzip.NewWriter(file)
				defer writer.Close()

				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					return err
				}

				_, err = writer.Write(body)
				if err != nil {
					return err
				}
			} else {
				_, err = io.Copy(file, res.Body)
				if err != nil {
					return err
				}
			}
		} else {
			if format == "mmdb" {
				_, err = io.Copy(file, res.Body)
				if err != nil {
					return err
				}
			} else {
				err := unzipWrite(file, res.Body)
				if err != nil {
					return err
				}
			}
		}

		fmt.Printf("Database %s saved successfully.", name)
	}

	return nil
}

func unzipWrite(file *os.File, data io.Reader) error {
	unzipData, err := gzip.NewReader(data)
	if err != nil {
		return err
	}
	defer unzipData.Close()

	_, err = io.Copy(file, unzipData)
	if err != nil {
		return err
	}

	return nil
}
