package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jszwec/csvutil"
)

func outputJSON(d interface{}) error {
	jsonEnc := json.NewEncoder(os.Stdout)
	jsonEnc.SetIndent("", "  ")
	return jsonEnc.Encode(d)
}

func outputCSV(v interface{}) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)

	if err := csvEnc.Encode(v); err != nil {
		return err
	}
	csvWriter.Flush()

	return nil
}

func outputCSVBatchCore(core ipinfo.BatchCore) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)

	// print entries.
	for _, v := range core {
		if err := csvEnc.Encode(v); err != nil {
			return err
		}
		csvWriter.Flush()
	}

	return nil
}

func outputFriendlyCore(d *ipinfo.Core) {
	var printline func(name string, val string)

	fmtHdr := color.New(color.Bold, color.FgWhite)
	fmtEntry := color.New(color.FgCyan)
	fmtVal := color.New(color.FgGreen)

	printlineGen := func(entryLen string) func(string, string) {
		return func(name string, val string) {
			fmt.Printf(
				"- %s %s\n",
				fmtEntry.Sprintf("%-"+entryLen+"s", name),
				fmtVal.Sprintf("%v", val),
			)
		}
	}

	fmtHdr.Println("Core")
	if d.Bogon {
		printline = printlineGen("5")
	} else {
		printline = printlineGen("12")
	}
	printline("IP", d.IP.String())
	if d.Bogon {
		// exit early after printing bogon status.
		printline("Bogon", fmt.Sprintf("%v", d.Bogon))
		return
	}
	printline("Anycast", fmt.Sprintf("%v", d.Anycast))
	printline("Hostname", d.Hostname)
	printline("City", d.City)
	printline("Region", d.Region)
	printline("Country", fmt.Sprintf("%v (%v)", d.CountryName, d.Country))
	printline("Location", d.Location)
	printline("Organization", d.Org)
	printline("Postal", d.Postal)
	printline("Timezone", d.Timezone)
	if d.ASN != nil {
		fmt.Println()
		fmtHdr.Println("ASN")
		printline = printlineGen("6")
		printline("ID", d.ASN.ASN)
		printline("Name", d.ASN.Name)
		printline("Domain", d.ASN.Domain)
		printline("Route", d.ASN.Route)
		printline("Type", d.ASN.Type)
	}
	if d.Company != nil {
		fmt.Println()
		fmtHdr.Println("Company")
		printline = printlineGen("6")
		printline("Name", d.Company.Name)
		printline("Domain", d.Company.Domain)
		printline("Type", d.Company.Type)
	}
	if d.Carrier != nil {
		fmt.Println()
		fmtHdr.Println("Carrier")
		printline = printlineGen("4")
		printline("Name", d.Carrier.Name)
		printline("MCC", d.Carrier.MCC)
		printline("MNC", d.Carrier.MNC)
	}
	if d.Privacy != nil {
		fmt.Println()
		fmtHdr.Println("Privacy")
		printline = printlineGen("7")
		printline("VPN", fmt.Sprintf("%v", d.Privacy.VPN))
		printline("Proxy", fmt.Sprintf("%v", d.Privacy.Proxy))
		printline("Tor", fmt.Sprintf("%v", d.Privacy.Tor))
		printline("Hosting", fmt.Sprintf("%v", d.Privacy.Hosting))
	}
	if d.Abuse != nil {
		fmt.Println()
		fmtHdr.Println("Abuse")
		printline = printlineGen("7")
		printline("Address", d.Abuse.Address)
		printline("Country", fmt.Sprintf("%v (%v)", d.Abuse.CountryName, d.Abuse.Country))
		printline("Email", d.Abuse.Email)
		printline("Name", d.Abuse.Name)
		printline("Network", d.Abuse.Network)
		printline("Phone", d.Abuse.Phone)
	}
	if d.Domains != nil && d.Domains.Total > 0 {
		fmt.Println()
		fmtHdr.Println("Domains")
		printline = printlineGen("8")
		printline("Total", fmt.Sprintf("%v", d.Domains.Total))
		if len(d.Domains.Domains) > 0 {
			printline("Examples", d.Domains.Domains[0])
			if len(d.Domains.Domains) > 1 {
				for _, d := range d.Domains.Domains[1:] {
					fmt.Printf(
						"           %v\n", fmtVal.Sprintf("%v", d),
					)
				}
			}
		}
	}
}

func outputFieldBatchCore(
	core ipinfo.BatchCore,
	field string,
	header bool,
	fieldOnly bool,
) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)
	csvEnc.AutoHeader = false

	// TODO the dread of not having macros... we can simplify code length here
	// with reflection but until then this will have to do.
	switch field {
	case "ip":
		if header {
			fmt.Printf("ip\n")
		}

		for _, d := range core {
			fmt.Printf("%v\n", d.IP)
		}
	case "hostname":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("hostname\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Hostname)
		}
	case "anycast":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("anycast\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Anycast)
		}
	case "city":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("city\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.City)
		}
	case "region":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("region\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Region)
		}
	case "country":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("country\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Country)
		}
	case "country_name":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("country_name\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.CountryName)
		}
	case "loc":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("loc\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Location)
		}
	case "org":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("org\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Org)
		}
	case "postal":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("postal\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Postal)
		}
	case "timezone":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("timezone\n")
		}

		for _, d := range core {
			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Timezone)
		}
	case "asn":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CoreASN{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.ASN); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "asn.id":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("asn_id\n")
		}

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.ASN.ASN)
		}
	case "asn.name", "asn.asn":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("asn_name\n")
		}

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.ASN.Name)
		}
	case "asn.domain":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("asn_domain\n")
		}

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.ASN.Domain)
		}
	case "asn.route":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("asn_route\n")
		}

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.ASN.Route)
		}
	case "asn.type":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("asn_type\n")
		}

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.ASN.Type)
		}
	case "company":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CoreCompany{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Company == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.Company); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "company.name":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("company_name\n")
		}

		for _, d := range core {
			if d.Company == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Company.Name)
		}
	case "company.domain":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("company_domain\n")
		}

		for _, d := range core {
			if d.Company == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Company.Domain)
		}
	case "company.type":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("company_type\n")
		}

		for _, d := range core {
			if d.Company == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Company.Type)
		}
	case "carrier":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CoreCarrier{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.Carrier); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "carrier.name":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("carrier_name\n")
		}

		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Carrier.Name)
		}
	case "carrier.mcc":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("carrier_mcc\n")
		}

		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Carrier.MCC)
		}
	case "carrier.mnc":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("carrier_mnc\n")
		}

		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Carrier.MNC)
		}
	case "privacy":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CorePrivacy{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.Privacy); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "privacy.vpn":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("privacy_vpn\n")
		}

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Privacy.VPN)
		}
	case "privacy.proxy":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("privacy_proxy\n")
		}

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Privacy.Proxy)
		}
	case "privacy.tor":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("privacy_tor\n")
		}

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Privacy.Tor)
		}
	case "privacy.hosting":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("privacy_hosting\n")
		}

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Privacy.Hosting)
		}
	case "abuse":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CoreAbuse{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.Abuse); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "abuse.address":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_address\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.Address)
		}
	case "abuse.country":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_country\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%s%v\"\n", d.Abuse.Country)
		}
	case "abuse.country_name":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_country_name\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.CountryName)
		}
	case "abuse.email":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_email\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.Email)
		}
	case "abuse.name":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_name\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.Name)
		}
	case "abuse.network":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_network\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.Network)
		}
	case "abuse.phone":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("abuse_phone\n")
		}

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Abuse.Phone)
		}
	case "domains":
		if !fieldOnly {
			fmt.Printf("ip,")
		}
		if err := csvEnc.EncodeHeader(ipinfo.CoreDomains{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Domains == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			if err := csvEnc.Encode(d.Domains); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "domains.total":
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("domains_total\n")
		}

		for _, d := range core {
			if d.Domains == nil {
				continue
			}

			if !fieldOnly {
				fmt.Printf("%s,", d.IP)
			}
			fmt.Printf("%v\n", d.Domains.Total)
		}
	default:
		if header {
			if !fieldOnly {
				fmt.Printf("ip,")
			}
			fmt.Printf("%s\n", field)
		}
	}

	return nil
}

func outputFieldBatchASNDetails(
	asnDetails ipinfo.BatchASNDetails,
	field string,
	header bool,
	fieldOnly bool,
) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)
	csvEnc.AutoHeader = false

	switch field {
	case "name":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("name\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Name)
		}
	case "country":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("country\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Country)
		}
	case "country_name":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("country_name\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.CountryName)
		}
	case "allocated":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("allocated\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Allocated)
		}
	case "registry":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("registry\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Registry)
		}
	case "domain":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("domain\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Domain)
		}
	case "num_ips":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("num_ips\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.NumIPs)
		}
	case "prefixes":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("prefixes\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Prefixes)
		}
	case "prefixes6":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("prefixes6\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Prefixes6)
		}
	case "peers":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("peers\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Peers)
		}
	case "upstreams":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("upstreams\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Upstreams)
		}
	case "downstreams":
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("downstreams\n")
		}

		for _, d := range asnDetails {
			if !fieldOnly {
				fmt.Printf("%s,", d.ASN)
			}
			fmt.Printf("%v\n", d.Downstreams)
		}
	default:
		if header {
			if !fieldOnly {
				fmt.Printf("asn,")
			}
			fmt.Printf("%s\n", field)
		}
	}

	return nil
}
