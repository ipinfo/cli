package main

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/spf13/pflag"
)

func printHelpSum() {
	fmt.Printf(
		`Usage: %s summarize [<opts>] <paths or '-' or cidrs or ip-range>

Description:
  Accepts file paths, '-' for stdin, CIDRs and IP ranges.

  # Lookup all IPs from stdin ('-' can be implied).
  $ %[1]s prips 8.8.8.0/24 | %[1]s summarize
  $ %[1]s prips 8.8.8.0/24 | %[1]s summarize -

  # Lookup all IPs in 2 files.
  $ %[1]s summarize /path/to/iplist1.txt /path/to/iplist2.txt

  # Lookup all IPs from CIDR.
  $ %[1]s summarize 8.8.8.0/24

  # Lookup all IPs from multiple CIDRs.
  $ %[1]s summarize 8.8.8.0/24 8.8.2.0/24 8.8.1.0/24

  # Lookup all IPs in an IP range.
  $ %[1]s summarize 8.8.8.0 8.8.8.255

Options:
  General:
    --token <tok>, -t <tok>
      use <tok> as API token.
    --help, -h
      show help.

  Formats:
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

	ips, err = getInputIPs(pflag.Args()[1:])
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		fmt.Println("no input ips")
		return nil
	}

	ii = prepareIpinfoClient(fTok)
	d, err := ii.GetIPSummary(ips)
	if err != nil {
		return err
	}

	if fJSON {
		return outputJSON(d)
	}

	// print pretty.
	var entryLen string
	var tmpEntryLen int
	header := color.New(color.Bold, color.FgWhite)
	entry := color.New(color.FgCyan)
	num := color.New(color.FgGreen)

	header.Println("Summary")
	headerPrint := func(name string, val uint64) {
		fmt.Printf(
			"- %s %s\n",
			entry.Sprintf("%-7s", name),
			num.Sprintf("%v", val),
		)
	}
	headerPrint("Total", d.Total)
	headerPrint("Unique", d.Unique)
	headerPrint("Anycast", d.Anycast)
	headerPrint("Bogon", d.Bogon)
	headerPrint("Mobile", d.Mobile)
	headerPrint("VPN", d.Privacy.VPN)
	headerPrint("Proxy", d.Privacy.Proxy)
	headerPrint("Hosting", d.Privacy.Hosting)
	headerPrint("Tor", d.Privacy.Tor)
	fmt.Println()

	header.Println("Top ASNs")
	topASNs := orderSummaryMapping(d.ASNs)
	entryLen = strconv.Itoa(longestKeyLen(topASNs))
	for _, asnSum := range topASNs {
		k := asnSum.k
		v := asnSum.v
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf("%-"+entryLen+"s", k),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}
	fmt.Println()

	header.Println("Top Usage Types")
	topUsageTypes := orderSummaryMapping(d.IPTypes)
	entryLen = strconv.Itoa(longestKeyLen(topUsageTypes))
	for _, usageTypeSum := range topUsageTypes {
		k := usageTypeSum.k
		if k == "isp" {
			k = "ISP"
		} else {
			k = strings.Title(k)
		}
		v := usageTypeSum.v
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf("%-"+entryLen+"s", k),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}
	fmt.Println()

	header.Println("Top Routes")
	topRoutes := orderSummaryMapping(d.Routes)
	entryLen = strconv.Itoa(longestKeyLen(topRoutes) + 2)
	for _, routesSum := range topRoutes {
		k := routesSum.k
		v := routesSum.v
		routeParts := strings.SplitN(k, " ", 2)
		asn := routeParts[0]
		route := routeParts[1]
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf(
				"%-"+entryLen+"s",
				fmt.Sprintf("%s (%s)", route, asn),
			),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}
	fmt.Println()

	header.Println("Top Countries")
	topCountries := orderSummaryMapping(d.Countries)
	for _, p := range topCountries {
		l := ipinfo.GetCountryName(p.k)
		if len(l) > tmpEntryLen {
			tmpEntryLen = len(l)
		}
	}
	entryLen = strconv.Itoa(tmpEntryLen)
	for _, countriesSum := range topCountries {
		k := countriesSum.k
		v := countriesSum.v
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf("%-"+entryLen+"s", ipinfo.GetCountryName(k)),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}
	fmt.Println()

	header.Println("Top Cities")
	topCities := orderSummaryMapping(d.Cities)
	entryLen = strconv.Itoa(longestKeyLen(topCities))
	for _, citiesSum := range topCities {
		k := citiesSum.k
		v := citiesSum.v
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf("%-"+entryLen+"s", k),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}
	fmt.Println()

	header.Println("Top Regions")
	topRegions := orderSummaryMapping(d.Regions)
	entryLen = strconv.Itoa(longestKeyLen(topRegions))
	for _, regionsSum := range topRegions {
		k := regionsSum.k
		v := regionsSum.v
		pct := (float64(v) / float64(d.Unique)) * 100
		fmt.Printf(
			"- %v %v\n",
			entry.Sprintf("%-"+entryLen+"s", k),
			num.Sprintf("%v (%.1f%%)", v, pct),
		)
	}

	if len(d.Carriers) > 0 {
		fmt.Println()
		header.Println("Top Carriers")
		topCarriers := orderSummaryMapping(d.Carriers)
		entryLen = strconv.Itoa(longestKeyLen(topCarriers))
		for _, carriersSum := range topCarriers {
			k := carriersSum.k
			v := carriersSum.v
			pct := (float64(v) / float64(d.Mobile)) * 100
			fmt.Printf(
				"- %v %v\n",
				entry.Sprintf("%-"+entryLen+"s", k),
				num.Sprintf("%v (%.1f%%)", v, pct),
			)
		}
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

func longestKeyLen(s []sumPair) int {
	longestK := 0
	for _, p := range s {
		if len(p.k) > longestK {
			longestK = len(p.k)
		}
	}
	return longestK
}
