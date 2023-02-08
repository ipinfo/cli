package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/cli/lib/complete"
	"github.com/ipinfo/cli/lib/complete/predict"
	"github.com/spf13/pflag"
)

type QuotaRequests struct {
	Day       int `json:"day"`
	Month     int `json:"month"`
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
}

type QuotaLimit struct {
	Daily   int `json:"daily"`
	Monthly int `json:"monthly"`
}

type QuotaFeatures struct {
	Core          QuotaLimit `json:"core"`
	Asn           QuotaLimit `json:"asn"`
	BasicAsn      QuotaLimit `json:"basic_asn"`
	Privacy       QuotaLimit `json:"privacy"`
	Company       QuotaLimit `json:"company"`
	Carrier       QuotaLimit `json:"carrier"`
	Ranges        QuotaLimit `json:"ranges"`
	Abuse         QuotaLimit `json:"abuse"`
	HostedDomains QuotaLimit `json:"hosted_domains"`
	Hostio        QuotaLimit `json:"hostio"`
	Whois         QuotaLimit `json:"whois"`
}

type QuotaBody struct {
	Token    string        `json:"token"`
	Requests QuotaRequests `json:"requests"`
	Features QuotaFeatures `json:"features"`
}

var completionsQuota = &complete.Command{
	Flags: map[string]complete.Predictor{
		"-d":         predict.Nothing,
		"--detailed": predict.Nothing,
		"-h":         predict.Nothing,
		"--help":     predict.Nothing,
	},
}

func printHelpQuota() {
	fmt.Printf(
		`Usage: %s quota [<opts>]

Options:
  --detailed, -d
    show a detailed view of all available limits.
    default: false.
  --help, -h
    show help.
`, progBase)
}

func cmdQuota() error {
	var fDetailed bool
	pflag.BoolVarP(&fDetailed, "detailed", "d", false, "detail view.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.Parse()

	if fHelp {
		printHelpQuota()
		return nil
	}

	token := gConfig.Token
	if token == "" {
		return errors.New("please login first to check quota")
	}

	res, err := http.Get("https://ipinfo.io/me?token=" + token)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	quota := &QuotaBody{}
	if err := json.NewDecoder(res.Body).Decode(quota); err != nil {
		return err
	}
	printStats(fDetailed, quota)

	return nil
}

func getUsageBar(percentage int) string {
	usageBar := "[" + strings.Repeat("#", percentage) + strings.Repeat(" ", 100-percentage) + "]"
	return fmt.Sprintf("%s %d%%\n", usageBar, percentage)
}

func printStats(detailed bool, quota *QuotaBody) {
	// Calculate the percentage of remaining quota.
	percentage := int(math.Round(float64(quota.Requests.Month) / float64(quota.Requests.Limit) * 100.0))

	// Pretty Print.
	fmtHdr := color.New(color.Bold, color.FgWhite)
	fmtEntry := color.New(color.FgCyan)
	fmtVal := color.New(color.FgGreen)
	pprint := func(name string, val int) {
		fmt.Printf(
			"- %s %s\n",
			fmtEntry.Sprintf("%-"+"25"+"s", name),
			fmtVal.Sprintf("%v", val),
		)
	}

	// Print stats.
	fmtHdr.Println("Usage")
	pprint("Total Requests", quota.Requests.Limit)
	pprint("Remaining Requests", quota.Requests.Remaining)
	pprint("Requests Made Today", quota.Requests.Day)
	pprint("Requests Made This Month", quota.Requests.Month)
	fmtHdr.Printf("\n%d%% of total %d used\n", percentage, quota.Requests.Limit)
	bar := getUsageBar(percentage)
	fmtVal.Println(bar)

	// If `detailed` print all values.
	if detailed {
		fmtVal.Println("Available Limits")
		// ASN Details.
		if quota.Features.Asn.Daily > 0 || quota.Features.Asn.Monthly > 0 {
			fmtHdr.Println("Asn")
			pprint("Daily Limit", quota.Features.Asn.Daily)
			pprint("Monthly Limit", quota.Features.Asn.Monthly)
		}

		// Basic_ASN Details.
		if quota.Features.BasicAsn.Daily > 0 || quota.Features.BasicAsn.Monthly > 0 {
			fmtHdr.Println("Basic_ASN")
			pprint("Daily Limit", quota.Features.BasicAsn.Daily)
			pprint("Monthly Limit", quota.Features.BasicAsn.Monthly)
		}

		// Privacy Details.
		if quota.Features.Privacy.Daily > 0 || quota.Features.Privacy.Monthly > 0 {
			fmtHdr.Println("Privacy")
			pprint("Daily Limit", quota.Features.Privacy.Daily)
			pprint("Monthly Limit", quota.Features.Privacy.Monthly)
		}

		// Company Details.
		if quota.Features.Company.Daily > 0 || quota.Features.Company.Monthly > 0 {
			fmtHdr.Println("Company")
			pprint("Daily Limit", quota.Features.Company.Daily)
			pprint("Monthly Limit", quota.Features.Company.Monthly)
		}

		// Carrier Details.
		if quota.Features.Carrier.Daily > 0 || quota.Features.Carrier.Monthly > 0 {
			fmtHdr.Println("Carrier")
			pprint("Daily Limit", quota.Features.Carrier.Daily)
			pprint("Monthly Limit", quota.Features.Carrier.Monthly)
		}

		// Ranges Details.
		if quota.Features.Ranges.Daily > 0 || quota.Features.Ranges.Monthly > 0 {
			fmtHdr.Println("Ranges")
			pprint("Daily Limit", quota.Features.Ranges.Daily)
			pprint("Monthly Limit", quota.Features.Ranges.Monthly)
		}

		// Abuse Details.
		if quota.Features.Abuse.Daily > 0 || quota.Features.Abuse.Monthly > 0 {
			fmtHdr.Println("Abuse")
			pprint("Daily Limit", quota.Features.Abuse.Daily)
			pprint("Monthly Limit", quota.Features.Abuse.Monthly)
		}

		// Hosted_Domains Details.
		if quota.Features.HostedDomains.Daily > 0 || quota.Features.HostedDomains.Monthly > 0 {
			fmtHdr.Println("Hosted_Domains")
			pprint("Daily Limit", quota.Features.HostedDomains.Daily)
			pprint("Monthly Limit", quota.Features.HostedDomains.Monthly)
		}

		// Hostio Details.
		if quota.Features.Hostio.Daily > 0 || quota.Features.Hostio.Monthly > 0 {
			fmtHdr.Println("Hostio")
			pprint("Daily Limit", quota.Features.Hostio.Daily)
			pprint("Monthly Limit", quota.Features.Hostio.Monthly)
		}

		// WHOIS Details.
		if quota.Features.Whois.Daily > 0 || quota.Features.Whois.Monthly > 0 {
			fmtHdr.Println("WHOIS")
			pprint("Daily Limit", quota.Features.Whois.Daily)
			pprint("Monthly Limit", quota.Features.Whois.Monthly)
		}
	}

}
