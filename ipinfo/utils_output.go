package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jszwec/csvutil"
)

var coreFields = []string{
	"ip",
	"hostname",
	"anycast",
	"city",
	"region",
	"country",
	"country_name",
	"loc",
	"org",
	"postal",
	"timezone",
	"asn",
	"asn.id",
	"asn.name",
	"asn.asn",
	"asn.domain",
	"asn.route",
	"asn.type",
	"company",
	"company.name",
	"company.domain",
	"company.type",
	"carrier",
	"carrier.name",
	"carrier.mcc",
	"carrier.mnc",
	"privacy",
	"privacy.vpn",
	"privacy.proxy",
	"privacy.tor",
	"privacy.hosting",
	"abuse",
	"abuse.address",
	"abuse.country",
	"abuse.country_name",
	"abuse.email",
	"abuse.name",
	"abuse.network",
	"abuse.phone",
	"domains",
	"domains.total",
}

var asnFields = []string{
	"id",
	"asn",
	"name",
	"country",
	"country_name",
	"allocated",
	"registry",
	"domain",
	"num_ips",
	"prefixes",
	"prefixes6",
	"peers",
	"upstreams",
	"downstreams",
}

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
	fields []string,
	header bool,
	inclIP bool,
) error {
	// error on bad field.
	for _, f := range fields {
		hasField := false
		for _, coreF := range coreFields {
			if coreF == f {
				hasField = true
				break
			}
		}
		if !hasField {
			errStr := "field '%v' is invalid; the following are allowed:"
			errStr += "  "+strings.Join(coreFields, "\n  ")
			return fmt.Errorf(errStr, f)
		}
	}

	// see if we should include IP as well.
	if header && inclIP {
		// only include it automatically if not already specified in list.
		hasIPField := false
		for _, f := range fields {
			if f == "ip" {
				hasIPField = true
				break
			}
		}
		if !hasIPField {
			fmt.Printf("ip,")
		}
	}

	hdrs := make([]string, 0, len(fields))
	rowFuncs := make([]func(*ipinfo.Core) string, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "asn":
			hdrs = append(
				hdrs,
				"asn_id",
				"asn_name",
				"asn_domain",
				"asn_route",
				"asn_type",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreASN)
		case "company":
			hdrs = append(
				hdrs,
				"company_name",
				"company_domain",
				"company_type",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreCompany)
		case "carrier":
			hdrs = append(
				hdrs,
				"carrier_name",
				"carrier_mcc",
				"carrier_mnc",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreCarrier)
		case "privacy":
			hdrs = append(
				hdrs,
				"privacy_vpn",
				"privacy_proxy",
				"privacy_tor",
				"privacy_hosting",
			)
			rowFuncs = append(rowFuncs, outputFieldCorePrivacy)
		case "abuse":
			hdrs = append(
				hdrs,
				"abuse_address",
				"abuse_country",
				"abuse_country_name",
				"abuse_email",
				"abuse_name",
				"abuse_network",
				"abuse_phone",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreAbuse)
		case "domains":
			hdrs = append(
				hdrs,
				"domains_total",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreDomains)
		default:
			hdrs = append(hdrs, strings.ReplaceAll(f, ".", "_"))
		}

		// funcs now.
		switch f {
		case "ip":
			rowFuncs = append(rowFuncs, outputFieldCoreIP)
		case "hostname":
			rowFuncs = append(rowFuncs, outputFieldCoreHostname)
		case "anycast":
			rowFuncs = append(rowFuncs, outputFieldCoreAnycast)
		case "city":
			rowFuncs = append(rowFuncs, outputFieldCoreCity)
		case "country":
			rowFuncs = append(rowFuncs, outputFieldCoreCountry)
		case "country_name":
			rowFuncs = append(rowFuncs, outputFieldCoreCountryName)
		case "loc":
			rowFuncs = append(rowFuncs, outputFieldCoreLoc)
		case "org":
			rowFuncs = append(rowFuncs, outputFieldCoreOrg)
		case "postal":
			rowFuncs = append(rowFuncs, outputFieldCorePostal)
		case "timezone":
			rowFuncs = append(rowFuncs, outputFieldCoreTimezone)
		case "asn.id":
			rowFuncs = append(rowFuncs, outputFieldCoreASNId)
		case "asn.name", "asn.asn":
			rowFuncs = append(rowFuncs, outputFieldCoreASNName)
		case "asn.domain":
			rowFuncs = append(rowFuncs, outputFieldCoreASNDomain)
		case "asn.route":
			rowFuncs = append(rowFuncs, outputFieldCoreASNRoute)
		case "asn.type":
			rowFuncs = append(rowFuncs, outputFieldCoreASNType)
		case "company.name":
			rowFuncs = append(rowFuncs, outputFieldCoreCompanyName)
		case "company.domain":
			rowFuncs = append(rowFuncs, outputFieldCoreCompanyDomain)
		case "company.type":
			rowFuncs = append(rowFuncs, outputFieldCoreCompanyType)
		case "carrier.name":
			rowFuncs = append(rowFuncs, outputFieldCoreCarrierName)
		case "carrier.mcc":
			rowFuncs = append(rowFuncs, outputFieldCoreCarrierMCC)
		case "carrier.mnc":
			rowFuncs = append(rowFuncs, outputFieldCoreCarrierMNC)
		case "privacy.vpn":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyVPN)
		case "privacy.proxy":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyProxy)
		case "privacy.tor":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyTor)
		case "privacy.hosting":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyHosting)
		case "abuse.address":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseAddress)
		case "abuse.country":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseCountry)
		case "abuse.country_name":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseCountryName)
		case "abuse.email":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseEmail)
		case "abuse.name":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseName)
		case "abuse.network":
			rowFuncs = append(rowFuncs, outputFieldCoreAbuseNetwork)
		case "abuse.phone":
			rowFuncs = append(rowFuncs, outputFieldCoreAbusePhone)
		case "domains.total":
			rowFuncs = append(rowFuncs, outputFieldCoreDomainsTotal)
		}
	}

	fmt.Println(strings.Join(hdrs, ","))
	for _, d := range core {
		row := make([]string, len(rowFuncs))
		for i, rowFunc := range rowFuncs {
			row[i] = rowFunc(d)
		}
		fmt.Println(strings.Join(row, ","))
	}

	return nil
}

func outputFieldBatchASNDetails(
	asnDetails ipinfo.BatchASNDetails,
	fields []string,
	header bool,
	inclASNId bool,
) error {
	// error on bad field.
	for _, f := range fields {
		hasField := false
		for _, asnF := range asnFields {
			if asnF == f {
				hasField = true
				break
			}
		}
		if !hasField {
			errStr := "field '%v' is invalid; the following are allowed:"
			errStr += "  "+strings.Join(asnFields, "\n  ")
			return fmt.Errorf(errStr, f)
		}
	}

	// see if we should include ASN name as well.
	if header && inclASNId {
		// only include it automatically if not already specified in list.
		hasASNIdField := false
		for _, f := range fields {
			if f == "id" || f == "asn" {
				hasASNIdField = true
				break
			}
		}
		if !hasASNIdField {
			fmt.Printf("id,")
		}
	}

	hdrs := make([]string, 0, len(fields))
	rowFuncs := make([]func(*ipinfo.ASNDetails) string, 0, len(fields))
	for _, f := range fields {
		hdrs = append(hdrs, strings.ReplaceAll(f, ".", "_"))

		switch f {
		case "id", "asn":
			rowFuncs = append(rowFuncs, outputFieldASNId)
		case "name":
			rowFuncs = append(rowFuncs, outputFieldASNName)
		case "country":
			rowFuncs = append(rowFuncs, outputFieldASNCountry)
		case "country_name":
			rowFuncs = append(rowFuncs, outputFieldASNCountryName)
		case "allocated":
			rowFuncs = append(rowFuncs, outputFieldASNAllocated)
		case "registry":
			rowFuncs = append(rowFuncs, outputFieldASNRegistry)
		case "domain":
			rowFuncs = append(rowFuncs, outputFieldASNDomain)
		case "num_ips":
			rowFuncs = append(rowFuncs, outputFieldASNNumIPs)
		case "prefixes":
			rowFuncs = append(rowFuncs, outputFieldASNPrefixes)
		case "prefixes6":
			rowFuncs = append(rowFuncs, outputFieldASNPrefixes6)
		case "peers":
			rowFuncs = append(rowFuncs, outputFieldASNPeers)
		case "upstreams":
			rowFuncs = append(rowFuncs, outputFieldASNUpstreams)
		case "downstreams":
			rowFuncs = append(rowFuncs, outputFieldASNDownstreams)
		}
	}

	fmt.Println(strings.Join(hdrs, ","))
	for _, d := range asnDetails {
		row := make([]string, len(rowFuncs))
		for i, rowFunc := range rowFuncs {
			row[i] = rowFunc(d)
		}
		fmt.Println(strings.Join(row, ","))
	}

	return nil
}

func outputFieldCoreIP(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.IP)
}

func outputFieldCoreHostname(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Hostname)
}

func outputFieldCoreAnycast(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Anycast)
}

func outputFieldCoreCity(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.City)
}

func outputFieldCoreRegion(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Region)
}

func outputFieldCoreCountry(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Country)
}

func outputFieldCoreCountryName(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.CountryName)
}

func outputFieldCoreLoc(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Location)
}

func outputFieldCoreOrg(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Org)
}

func outputFieldCorePostal(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Postal)
}

func outputFieldCoreTimezone(core *ipinfo.Core) string {
	return fmt.Sprintf("%v", core.Timezone)
}

func outputFieldCoreASN(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ",,,,"
	}
	return fmt.Sprintf(
		"%v,%v,%v,%v,%v",
		core.ASN.ASN,
		core.ASN.Name,
		core.ASN.Domain,
		core.ASN.Route,
		core.ASN.Type,
	)
}

func outputFieldCoreASNId(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.ASN.ASN)
}

func outputFieldCoreASNName(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.ASN.Name)
}

func outputFieldCoreASNDomain(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.ASN.Domain)
}

func outputFieldCoreASNRoute(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.ASN.Route)
}

func outputFieldCoreASNType(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.ASN.Type)
}

func outputFieldCoreCompany(core *ipinfo.Core) string {
	if core.Company == nil {
		return ",,"
	}
	return fmt.Sprintf(
		"%v,%v,%v",
		core.Company.Name,
		core.Company.Domain,
		core.Company.Type,
	)
}

func outputFieldCoreCompanyName(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Company.Name)
}

func outputFieldCoreCompanyDomain(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Company.Domain)
}

func outputFieldCoreCompanyType(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Company.Type)
}

func outputFieldCoreCarrier(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ",,"
	}
	return fmt.Sprintf(
		"%v,%v,%v",
		core.Carrier.Name,
		core.Carrier.MCC,
		core.Carrier.MNC,
	)
}

func outputFieldCoreCarrierName(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Carrier.Name)
}

func outputFieldCoreCarrierMCC(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Carrier.MCC)
}

func outputFieldCoreCarrierMNC(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Carrier.MNC)
}

func outputFieldCorePrivacy(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ",,,"
	}
	return fmt.Sprintf(
		"%v,%v,%v,%v",
		core.Privacy.VPN,
		core.Privacy.Proxy,
		core.Privacy.Tor,
		core.Privacy.Hosting,
	)
}

func outputFieldCorePrivacyVPN(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Privacy.VPN)
}

func outputFieldCorePrivacyProxy(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Privacy.Proxy)
}

func outputFieldCorePrivacyTor(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Privacy.Tor)
}

func outputFieldCorePrivacyHosting(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Privacy.Hosting)
}

func outputFieldCoreAbuse(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ",,,,,,"
	}
	return fmt.Sprintf(
		"%v,%v,%v,%v,%v,%v,%v",
		core.Abuse.Address,
		core.Abuse.Country,
		core.Abuse.CountryName,
		core.Abuse.Email,
		core.Abuse.Name,
		core.Abuse.Network,
		core.Abuse.Phone,
	)
}

func outputFieldCoreAbuseAddress(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Address)
}

func outputFieldCoreAbuseCountry(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Country)
}

func outputFieldCoreAbuseCountryName(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.CountryName)
}

func outputFieldCoreAbuseEmail(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Email)
}

func outputFieldCoreAbuseName(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Name)
}

func outputFieldCoreAbuseNetwork(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Network)
}

func outputFieldCoreAbusePhone(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Abuse.Phone)
}

func outputFieldCoreDomains(core *ipinfo.Core) string {
	if core.Domains == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Domains.Total)
}

func outputFieldCoreDomainsTotal(core *ipinfo.Core) string {
	if core.Domains == nil {
		return ""
	}
	return fmt.Sprintf("%v", core.Domains.Total)
}

func outputFieldASNId(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.ASN)
}

func outputFieldASNName(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Name)
}

func outputFieldASNCountry(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Country)
}

func outputFieldASNCountryName(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.CountryName)
}

func outputFieldASNAllocated(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Allocated)
}

func outputFieldASNRegistry(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Registry)
}

func outputFieldASNDomain(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Domain)
}

func outputFieldASNNumIPs(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.NumIPs)
}

func outputFieldASNPrefixes(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Prefixes)
}

func outputFieldASNPrefixes6(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Prefixes6)
}

func outputFieldASNPeers(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Peers)
}

func outputFieldASNUpstreams(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Upstreams)
}

func outputFieldASNDownstreams(d *ipinfo.ASNDetails) string {
	return fmt.Sprintf("%v", d.Downstreams)
}
