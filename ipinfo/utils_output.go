package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jszwec/csvutil"
	"gopkg.in/yaml.v3"
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
	"privacy",
	"privacy.vpn",
	"privacy.proxy",
	"privacy.tor",
	"privacy.relay",
	"privacy.hosting",
	"privacy.service",
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

// Converts `i` to a single CSV encoded line, where `i` may be a `string`,
// `[]string`, `bool` or any value that can be converted to a string.
func encodeToCsvLine(i interface{}) string {
	var arr []string
	switch v := i.(type) {
	case []string:
		arr = v
	case string:
		arr = []string{v}
	case bool:
		arr = []string{strconv.FormatBool(v)}
	default:
		arr = []string{fmt.Sprintf("%v", v)}
	}

	// writing to string buffer.
	buf := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buf)
	if err := csvWriter.Write(arr); err != nil {
		panic(err)
	}
	csvWriter.Flush()
	encodedStr := buf.String()
	end := len(encodedStr) - 1
	if end < 0 {
		end = 0
	}
	return encodedStr[:end]
}

func outputYAML(d interface{}) error {
	yamlEnc := yaml.NewEncoder(os.Stdout)
	return yamlEnc.Encode(d)
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
	if d.Country == "" {
		printline("Country", "")
	} else {
		printline("Country", fmt.Sprintf("%v (%v)", d.CountryName, d.Country))
		printline("Currency", fmt.Sprintf("%v (%v)", d.CountryCurrency.Code, d.CountryCurrency.Symbol))
	}
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
	if d.Privacy != nil {
		fmt.Println()
		fmtHdr.Println("Privacy")
		printline = printlineGen("7")
		printline("VPN", fmt.Sprintf("%v", d.Privacy.VPN))
		printline("Proxy", fmt.Sprintf("%v", d.Privacy.Proxy))
		printline("Tor", fmt.Sprintf("%v", d.Privacy.Tor))
		printline("Relay", fmt.Sprintf("%v", d.Privacy.Relay))
		printline("Hosting", fmt.Sprintf("%v", d.Privacy.Hosting))
		printline("Service", fmt.Sprintf("%v", d.Privacy.Service))
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
			errStr += "  " + strings.Join(coreFields, "\n  ")
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
			fields = append([]string{"ip"}, fields...)
		}
	}

	hdrs := make([]string, 0, len(fields))
	rowFuncs := make([]func(*ipinfo.Core) string, 0, len(fields))
	var errs []error
	// We aim to print only a single error message, even if multiple
	// fields are related to a single piece of premium data.
	asnErrorHandled := false
	companyErrorHandled := false
	privacyErrorHandled := false
	abuseErrorHandled := false
	domainsErrorHandled := false
	for _, f := range fields {
		if strings.HasPrefix(f, "asn") && !asnErrorHandled {
			if err := checkPremiumData(core, "asn"); err != nil {
				errs = append(errs, err)
				asnErrorHandled = true
				continue
			}
		} else if strings.HasPrefix(f, "company") && !companyErrorHandled {
			if err := checkPremiumData(core, "company"); err != nil {
				errs = append(errs, err)
				companyErrorHandled = true
				continue
			}
		} else if strings.HasPrefix(f, "privacy") && !privacyErrorHandled {
			if err := checkPremiumData(core, "privacy"); err != nil {
				errs = append(errs, err)
				privacyErrorHandled = true
				continue
			}
		} else if strings.HasPrefix(f, "abuse") && !abuseErrorHandled {
			if err := checkPremiumData(core, "abuse"); err != nil {
				errs = append(errs, err)
				abuseErrorHandled = true
				continue
			}
		} else if strings.HasPrefix(f, "domains") && !domainsErrorHandled {
			if err := checkPremiumData(core, "domains"); err != nil {
				errs = append(errs, err)
				domainsErrorHandled = true
				continue
			}
		}

		switch f {
		case "asn":
			if asnErrorHandled {
				continue
			}
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
			if companyErrorHandled {
				continue
			}
			hdrs = append(
				hdrs,
				"company_name",
				"company_domain",
				"company_type",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreCompany)
		case "privacy":
			if privacyErrorHandled {
				continue
			}
			hdrs = append(
				hdrs,
				"privacy_vpn",
				"privacy_proxy",
				"privacy_tor",
				"privacy_relay",
				"privacy_hosting",
				"privacy_service",
			)
			rowFuncs = append(rowFuncs, outputFieldCorePrivacy)
		case "abuse":
			if abuseErrorHandled {
				continue
			}
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
			if domainsErrorHandled {
				continue
			}
			hdrs = append(
				hdrs,
				"domains_total",
			)
			rowFuncs = append(rowFuncs, outputFieldCoreDomains)
		default:
			if asnErrorHandled && strings.HasPrefix(f, "asn.") {
				continue
			} else if companyErrorHandled && strings.HasPrefix(f, "company.") {
				continue
			} else if privacyErrorHandled && strings.HasPrefix(f, "privacy.") {
				continue
			} else if abuseErrorHandled && strings.HasPrefix(f, "abuse.") {
				continue
			} else if domainsErrorHandled && strings.HasPrefix(f, "domains.") {
				continue
			} else {
				hdrs = append(hdrs, strings.ReplaceAll(f, ".", "_"))
			}

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
		case "region":
			rowFuncs = append(rowFuncs, outputFieldCoreRegion)
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
		case "privacy.vpn":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyVPN)
		case "privacy.proxy":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyProxy)
		case "privacy.tor":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyTor)
		case "privacy.relay":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyRelay)
		case "privacy.hosting":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyHosting)
		case "privacy.service":
			rowFuncs = append(rowFuncs, outputFieldCorePrivacyService)
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

	if len(hdrs) > 0 {
		fmt.Println(strings.Join(hdrs, ","))
		for _, d := range core {
			row := make([]string, len(rowFuncs))
			for i, rowFunc := range rowFuncs {
				row[i] = rowFunc(d)
			}
			fmt.Println(strings.Join(row, ","))
		}
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
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
			errStr += "  " + strings.Join(asnFields, "\n  ")
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
			fields = append([]string{"id"}, fields...)
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

// The core struct within the BatchCore type has 5 types of premium data.
// 1. ASN 2. Company 3. Privacy 4. Abuse 5. Domains.
// These 5 fields will be populated depending upon the type of the token
// through which the API call is made.
// If the account behind the token is on a free plan, all 5 premium data fields will be nil.
// If the account behind the token is on a basic plan, only the ASN field will get populated
// and the rest would be nil.
// If the account behind the token is on a standard plan, only ASN and Privacy will be populated.
// If the account behind the token is on a business plan, all 5 fields will get populated.
//
// The aim is to determine the user's token type and check what type of data the user prompted
// the CLI to output.
func checkPremiumData(core ipinfo.BatchCore, dataType string) error {
	for _, coreEntry := range core {
		if coreEntry.ASN == nil && coreEntry.Company == nil && coreEntry.Privacy == nil && coreEntry.Abuse == nil && coreEntry.Domains == nil {
			return fmt.Errorf("this token doesn't have permissions to access %s data", dataType)
		} else if coreEntry.ASN != nil && coreEntry.Company == nil && coreEntry.Privacy == nil && coreEntry.Abuse == nil && coreEntry.Domains == nil {
			if dataType == "company" || dataType == "privacy" || dataType == "abuse" || dataType == "domains" {
				return fmt.Errorf("this token doesn't have permissions to access %s data", dataType)
			}
			return nil
		} else if coreEntry.ASN != nil && coreEntry.Company == nil && coreEntry.Privacy != nil && coreEntry.Abuse == nil && coreEntry.Domains == nil {
			if dataType == "company" || dataType == "abuse" || dataType == "domains" {
				return fmt.Errorf("this token doesn't have permissions to access %s data", dataType)
			}
			return nil
		} else {
			return nil
		}
	}

	return nil
}

func outputFieldCoreIP(core *ipinfo.Core) string {
	return encodeToCsvLine(core.IP)
}

func outputFieldCoreHostname(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Hostname)
}

func outputFieldCoreAnycast(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Anycast)
}

func outputFieldCoreCity(core *ipinfo.Core) string {
	return encodeToCsvLine(core.City)
}

func outputFieldCoreRegion(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Region)
}

func outputFieldCoreCountry(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Country)
}

func outputFieldCoreCountryName(core *ipinfo.Core) string {
	return encodeToCsvLine(core.CountryName)
}

func outputFieldCoreLoc(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Location)
}

func outputFieldCoreOrg(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Org)
}

func outputFieldCorePostal(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Postal)
}

func outputFieldCoreTimezone(core *ipinfo.Core) string {
	return encodeToCsvLine(core.Timezone)
}

func outputFieldCoreASN(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ",,,,"
	}
	return encodeToCsvLine([]string{
		core.ASN.ASN,
		core.ASN.Name,
		core.ASN.Domain,
		core.ASN.Route,
		core.ASN.Type,
	})
}

func outputFieldCoreASNId(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return encodeToCsvLine(core.ASN.ASN)
}

func outputFieldCoreASNName(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return encodeToCsvLine(core.ASN.Name)
}

func outputFieldCoreASNDomain(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return encodeToCsvLine(core.ASN.Domain)
}

func outputFieldCoreASNRoute(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return encodeToCsvLine(core.ASN.Route)
}

func outputFieldCoreASNType(core *ipinfo.Core) string {
	if core.ASN == nil {
		return ""
	}
	return encodeToCsvLine(core.ASN.Type)
}

func outputFieldCoreCompany(core *ipinfo.Core) string {
	if core.Company == nil {
		return ",,"
	}
	return encodeToCsvLine([]string{
		core.Company.Name,
		core.Company.Domain,
		core.Company.Type,
	})
}

func outputFieldCoreCompanyName(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return encodeToCsvLine(core.Company.Name)
}

func outputFieldCoreCompanyDomain(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return encodeToCsvLine(core.Company.Domain)
}

func outputFieldCoreCompanyType(core *ipinfo.Core) string {
	if core.Company == nil {
		return ""
	}
	return encodeToCsvLine(core.Company.Type)
}

func outputFieldCorePrivacy(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ",,,,,"
	}
	return encodeToCsvLine([]string{
		strconv.FormatBool(core.Privacy.VPN),
		strconv.FormatBool(core.Privacy.Proxy),
		strconv.FormatBool(core.Privacy.Tor),
		strconv.FormatBool(core.Privacy.Relay),
		strconv.FormatBool(core.Privacy.Hosting),
		core.Privacy.Service,
	})
}

func outputFieldCorePrivacyVPN(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.VPN)
}

func outputFieldCorePrivacyProxy(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.Proxy)
}

func outputFieldCorePrivacyTor(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.Tor)
}

func outputFieldCorePrivacyRelay(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.Relay)
}

func outputFieldCorePrivacyHosting(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.Hosting)
}

func outputFieldCorePrivacyService(core *ipinfo.Core) string {
	if core.Privacy == nil {
		return ""
	}
	return encodeToCsvLine(core.Privacy.Service)
}

func outputFieldCoreAbuse(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ",,,,,,"
	}
	return encodeToCsvLine([]string{
		core.Abuse.Address,
		core.Abuse.Country,
		core.Abuse.CountryName,
		core.Abuse.Email,
		core.Abuse.Name,
		core.Abuse.Network,
		core.Abuse.Phone,
	})
}

func outputFieldCoreAbuseAddress(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Address)
}

func outputFieldCoreAbuseCountry(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Country)
}

func outputFieldCoreAbuseCountryName(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.CountryName)
}

func outputFieldCoreAbuseEmail(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Email)
}

func outputFieldCoreAbuseName(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Name)
}

func outputFieldCoreAbuseNetwork(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Network)
}

func outputFieldCoreAbusePhone(core *ipinfo.Core) string {
	if core.Abuse == nil {
		return ""
	}
	return encodeToCsvLine(core.Abuse.Phone)
}

func outputFieldCoreDomains(core *ipinfo.Core) string {
	if core.Domains == nil {
		return ""
	}
	return encodeToCsvLine(core.Domains.Total)
}

func outputFieldCoreDomainsTotal(core *ipinfo.Core) string {
	if core.Domains == nil {
		return ""
	}
	return encodeToCsvLine(core.Domains.Total)
}

func outputFieldASNId(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.ASN)
}

func outputFieldASNName(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Name)
}

func outputFieldASNCountry(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Country)
}

func outputFieldASNCountryName(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.CountryName)
}

func outputFieldASNAllocated(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Allocated)
}

func outputFieldASNRegistry(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Registry)
}

func outputFieldASNDomain(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Domain)
}

func outputFieldASNNumIPs(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.NumIPs)
}

func outputFieldASNPrefixes(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Prefixes)
}

func outputFieldASNPrefixes6(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Prefixes6)
}

func outputFieldASNPeers(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Peers)
}

func outputFieldASNUpstreams(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Upstreams)
}

func outputFieldASNDownstreams(d *ipinfo.ASNDetails) string {
	return encodeToCsvLine(d.Downstreams)
}
