package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

func printHelpSum() {
	fmt.Printf(
		`Usage: %s sum [<opts>] <paths or '-' or cidrs or ip-range>

Description:
  Accepts file paths, '-' for stdin, CIDRs and IP ranges.

  # Lookup all IPs from stdin ('-' can be implied).
  $ %[1]s prips 8.8.8.0/24 | %[1]s sum
  $ %[1]s prips 8.8.8.0/24 | %[1]s sum -

  # Lookup all IPs in 2 files.
  $ %[1]s sum /path/to/iplist1.txt /path/to/iplist2.txt

  # Lookup all IPs from CIDR.
  $ %[1]s sum 8.8.8.0/24

  # Lookup all IPs from multiple CIDRs.
  $ %[1]s sum 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # Lookup all IPs in an IP range.
  $ %[1]s sum 8.8.8.0 8.8.8.255

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Outputs:
    --pretty, -p
      output pretty format. (default)
    --json, -j
      output JSON format.
`, progBase)
}

func cmdSum() (err error) {
	var ips []net.IP
	var fTok string
	var fHelp bool
	var fPretty bool
	var fJSON bool

	pflag.StringVarP(&fTok, "token", "t", "", "the token to use.")
	pflag.BoolVarP(&fHelp, "help", "h", false, "show help.")
	pflag.BoolVarP(&fPretty, "pretty", "p", true, "output pretty format. (default)")
	pflag.BoolVarP(&fJSON, "json", "j", false, "output JSON format.")
	pflag.Parse()

	if fHelp {
		printHelpSum()
		return nil
	}

	if err := prepareIpinfoClient(fTok); err != nil {
		return err
	}

	args := pflag.Args()[1:]

	// check for stdin, implied or explicit.
	if len(args) == 0 || (len(args) == 1 && args[0] == "-") {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Println("** manual input mode **")
			fmt.Println("Enter all IPs, one per line:")
		}
		ips = ipsFromStdin()

		goto lookup
	}

	// check for IP range.
	if isIP(args[0]) {
		if len(args) != 2 {
			return errIPRangeRequiresTwoIPs
		}
		if !isIP(args[1]) {
			return errNotIP
		}

		ips, err = ipsFromRange(args[0], args[1])
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all CIDRs.
	if isCIDR(args[0]) {
		for _, arg := range args[1:] {
			if !isCIDR(arg) {
				return errNotCIDR
			}
		}

		ips, err = ipsFromCIDRs(args)
		if err != nil {
			return err
		}

		goto lookup
	}

	// check for all filepaths.
	if fileExists(args[0]) {
		for _, arg := range args[1:] {
			if !fileExists(arg) {
				return errNotFile
			}
		}

		ips, err = ipsFromFiles(args)
		if err != nil {
			return err
		}

		goto lookup
	}

lookup:

	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	d, err := ii.GetIPSummary(ips)
	if err != nil {
		return err
	}

	if fJSON {
		return outputJSON(d)
	}

	// print pretty.
	header := color.New(color.Bold, color.BgWhite, color.FgHiMagenta)

	header.Printf("                SUMMARY               ")
	fmt.Println()
	fmt.Printf("Total                          %v\n", d.Total)
	fmt.Printf("Unique                         %v\n", d.Unique)
	fmt.Printf("Anycast                        %v\n", d.Anycast)
	fmt.Printf("Bogon                          %v\n", d.Bogon)
	fmt.Printf("VPN                            %v\n", d.Privacy.VPN)
	fmt.Printf("Proxy                          %v\n", d.Privacy.Proxy)
	fmt.Printf("Hosting                        %v\n", d.Privacy.Hosting)
	fmt.Printf("Tor                            %v\n", d.Privacy.Tor)
	fmt.Println()
	header.Printf("                TOP ASNs              ")
	fmt.Println()
	topASNs := orderSummaryMapping(d.ASNs)
	for i, asnSum := range topASNs {
		k := asnSum.k
		v := asnSum.v

		asnParts := strings.SplitN(k, " ", 2)
		id := asnParts[0]
		name := asnParts[1]

		pct := (float64(v) / float64(d.Unique)) * 100
		barCnt := int(pct / 5)
		bar := createBarString(barCnt, 30)

		fmt.Printf(
			"%v. %-18v %13s\n",
			i+1, id, fmt.Sprintf("%v (%.1f%%)", v, pct),
		)
		fmt.Printf("   %v\n", name)
		fmt.Printf("   %-30s\n", bar)
		fmt.Println()
	}
	header.Printf("             TOP COUNTRIES            ")
	fmt.Println()
	topCountries := orderSummaryMapping(d.Countries)
	for i, countriesSum := range topCountries {
		k := countriesSum.k
		v := countriesSum.v

		pct := (float64(v) / float64(d.Unique)) * 100
		barCnt := int(pct / 5)
		bar := createBarString(barCnt, 30)

		fmt.Printf(
			"%v. %-18v %13s\n",
			i+1,
			ipinfo.GetCountryName(k),
			fmt.Sprintf("%v (%.1f%%)", v, pct),
		)
		fmt.Printf("   %-30s\n", bar)
		fmt.Println()
	}
	header.Printf("            TOP USAGE TYPES           ")
	fmt.Println()
	topUsageTypes := orderSummaryMapping(d.IPTypes)
	for i, usageTypeSum := range topUsageTypes {
		k := usageTypeSum.k
		if k == "isp" {
			k = "ISP"
		} else {
			k = strings.Title(k)
		}
		v := usageTypeSum.v

		pct := (float64(v) / float64(d.Unique)) * 100
		barCnt := int(pct / 5)
		bar := createBarString(barCnt, 30)

		fmt.Printf(
			"%v. %-18v %13s\n",
			i+1,
			k,
			fmt.Sprintf("%v (%.1f%%)", v, pct),
		)
		fmt.Printf("   %-30s\n", bar)
		fmt.Println()
	}

	return nil
}

/*
small utility for properly sorting summary results.

this is only needed because Golang maps don't guarantee ordered traversals.
when we decode from the raw JSON, which *is* sorted already, we lose that
sort order so have to regain it here.
*/

type sumPair struct {
	k string
	v uint64
}

type sumPairList []sumPair

func (s sumPairList) Len() int {
	return len(s)
}

func (s sumPairList) Less(i, j int) bool {
	return s[i].v < s[j].v
}

func (s sumPairList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func orderSummaryMapping(m map[string]uint64) []sumPair {
	s := make(sumPairList, 0, len(m))
	for k, v := range m {
		s = append(s, sumPair{k, v})
	}

	sort.Sort(sort.Reverse(s))
	return s
}
