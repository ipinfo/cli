package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
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
	"carrier",
	"carrier.name",
	"carrier.mcc",
	"carrier.mnc",
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

type OutputOption func(*outputOptions)

type outputOptions struct {
	accessibleData map[string]bool
	permissionsErr *PermissionsError
}

// Option to set accessible data
func WithAccessibleData(data map[string]bool) OutputOption {
	return func(opts *outputOptions) {
		opts.accessibleData = data
	}
}

// Option to set permissions error
func WithPermissionsError(err *PermissionsError) OutputOption {
	return func(opts *outputOptions) {
		opts.permissionsErr = err
	}
}

func outputFieldBatchCore(
	core ipinfo.BatchCore,
	fields []string,
	header bool,
	inclIP bool,
	opts ...OutputOption,
) error {
	options := outputOptions{}
	for _, opt := range opts {
		opt(&options)
	}
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
	for _, f := range fields {
		if options.accessibleData != nil && !options.accessibleData[f] {
			continue
		}
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
				"privacy_relay",
				"privacy_hosting",
				"privacy_service",
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

	if options.permissionsErr != nil {
		return fmt.Errorf(options.permissionsErr.Error())
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

type PermissionsError struct {
	MissingPermissions []string
	CarrierMessage     string
}

func (e *PermissionsError) Error() string {
	if len(e.MissingPermissions) > 0 {
		return fmt.Sprintf("This token doesn't have permissions to access the following data: %s", strings.Join(e.MissingPermissions, ", "))
	} else if e.CarrierMessage != "" {
		return e.CarrierMessage
	}

	return ""
}

func checkTokenPermissions(data *ipinfo.Core) (map[string]bool, *PermissionsError) {
	accessible := make(map[string]bool)
	var errInfo PermissionsError

	// Check ASN data access
	if data.ASN != nil {
		accessible["asn"] = true
	} else {
		errInfo.MissingPermissions = append(errInfo.MissingPermissions, "asn")
	}

	// Check Company data access
	if data.Company != nil {
		accessible["company"] = true
	} else {
		errInfo.MissingPermissions = append(errInfo.MissingPermissions, "company")
	}

	// Check Carrier data access
	if data.Carrier != nil {
		accessible["carrier"] = true
	} else {
		if data.ASN != nil && data.Company != nil && data.Privacy != nil && data.Abuse != nil && data.Domains != nil {
			// The highest plan behind the token is assumed on the presence of other data
			errInfo.CarrierMessage = "No carrier data associated with this IP address."
		} else {
			errInfo.MissingPermissions = append(errInfo.MissingPermissions, "carrier")
		}

	}

	// Check Privacy data access
	if data.Privacy != nil {
		accessible["privacy"] = true
	} else {
		errInfo.MissingPermissions = append(errInfo.MissingPermissions, "privacy")
	}

	// Check Abuse data access
	if data.Abuse != nil {
		accessible["abuse"] = true
	} else {
		errInfo.MissingPermissions = append(errInfo.MissingPermissions, "abuse")
	}

	// Check Domains data access
	if data.Domains != nil {
		accessible["domains"] = true
	} else {
		errInfo.MissingPermissions = append(errInfo.MissingPermissions, "domains")
	}

	if len(errInfo.MissingPermissions) == 0 && errInfo.CarrierMessage == "" {
		return accessible, nil
	}

	return accessible, &errInfo
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

func outputFieldCoreCarrier(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ",,"
	}
	return encodeToCsvLine([]string{
		core.Carrier.Name,
		core.Carrier.MCC,
		core.Carrier.MNC,
	})
}

func outputFieldCoreCarrierName(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return encodeToCsvLine(core.Carrier.Name)
}

func outputFieldCoreCarrierMCC(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return encodeToCsvLine(core.Carrier.MCC)
}

func outputFieldCoreCarrierMNC(core *ipinfo.Core) string {
	if core.Carrier == nil {
		return ""
	}
	return encodeToCsvLine(core.Carrier.MNC)
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
