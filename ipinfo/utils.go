package main

import (
	"bufio"
	"encoding/binary"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/jszwec/csvutil"
)

func isIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
}

func isCIDR(cidrStr string) bool {
	_, _, err := net.ParseCIDR(cidrStr)
	return err == nil
}

func isASN(asn string) bool {
	// check length.
	if len(asn) < 3 {
		return false
	}

	// ensure "AS" or "as" prefix.
	if !strings.HasPrefix(asn, "AS") && !strings.HasPrefix(asn, "as") {
		return false
	}

	// ensure number suffix.
	asnNumStr := asn[2:]
	if _, err := strconv.Atoi(asnNumStr); err != nil {
		return false
	}

	return true
}

func outputJSON(d interface{}) error {
	jsonEnc := json.NewEncoder(os.Stdout)
	jsonEnc.SetIndent("", "  ")
	return jsonEnc.Encode(d)
}

func outputFriendlyCore(d *ipinfo.Core) {
	header := color.New(color.Bold, color.BgWhite, color.FgHiMagenta)

	header.Printf("                 CORE                 ")
	fmt.Println()
	fmt.Printf("IP              %s\n", d.IP.String())
	fmt.Printf("Anycast         %v\n", d.Anycast)
	fmt.Printf("Hostname        %s\n", d.Hostname)
	fmt.Printf("City            %s\n", d.City)
	fmt.Printf("Region          %s\n", d.Region)
	fmt.Printf("Country         %s (%s)\n", d.CountryName, d.Country)
	fmt.Printf("Location        %s\n", d.Location)
	fmt.Printf("Organization    %s\n", d.Org)
	fmt.Printf("Postal          %s\n", d.Postal)
	fmt.Printf("Timezone        %s\n", d.Timezone)
	if d.ASN != nil {
		fmt.Println()
		header.Printf("                 ASN                  ")
		fmt.Println()
		fmt.Printf("ID              %s\n", d.ASN.ASN)
		fmt.Printf("Name            %s\n", d.ASN.Name)
		fmt.Printf("Domain          %s\n", d.ASN.Domain)
		fmt.Printf("Route           %s\n", d.ASN.Route)
		fmt.Printf("Type            %s\n", d.ASN.Type)
	}
	if d.Company != nil {
		fmt.Println()
		header.Printf("               COMPANY                ")
		fmt.Println()
		fmt.Printf("Name            %s\n", d.Company.Name)
		fmt.Printf("Domain          %s\n", d.Company.Domain)
		fmt.Printf("Type            %s\n", d.Company.Type)
	}
	if d.Carrier != nil {
		fmt.Println()
		header.Printf("               CARRIER                ")
		fmt.Println()
		fmt.Printf("Name            %s\n", d.Carrier.Name)
		fmt.Printf("MCC             %s\n", d.Carrier.MCC)
		fmt.Printf("MNC             %s\n", d.Carrier.MNC)
	}
	if d.Privacy != nil {
		fmt.Println()
		header.Printf("               PRIVACY                ")
		fmt.Println()
		fmt.Printf("VPN             %v\n", d.Privacy.VPN)
		fmt.Printf("Proxy           %v\n", d.Privacy.Proxy)
		fmt.Printf("Tor             %v\n", d.Privacy.Tor)
		fmt.Printf("Hosting         %v\n", d.Privacy.Hosting)
	}
	if d.Abuse != nil {
		fmt.Println()
		header.Printf("                ABUSE                 ")
		fmt.Println()
		fmt.Printf("Address         %s\n", d.Abuse.Address)
		fmt.Printf("Country         %s (%s)\n", d.Abuse.CountryName, d.Abuse.Country)
		fmt.Printf("Email           %s\n", d.Abuse.Email)
		fmt.Printf("Name            %s\n", d.Abuse.Name)
		fmt.Printf("Network         %s\n", d.Abuse.Network)
		fmt.Printf("Phone           %s\n", d.Abuse.Phone)
	}
	if d.Domains != nil && d.Domains.Total > 0 {
		fmt.Println()
		header.Printf("               DOMAINS                ")
		fmt.Println()
		fmt.Printf("Total           %v\n", d.Domains.Total)
		if len(d.Domains.Domains) > 0 {
			fmt.Printf("Examples     1: %s\n", d.Domains.Domains[0])
			if len(d.Domains.Domains) > 1 {
				for i, d := range d.Domains.Domains[1:] {
					fmt.Printf("             %v: %s\n", i+2, d)
				}
			}
		}
	}
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

func outputFieldBatchCore(core ipinfo.BatchCore, field string) error {
	csvWriter := csv.NewWriter(os.Stdout)
	csvEnc := csvutil.NewEncoder(csvWriter)
	csvEnc.AutoHeader = false

	// TODO the dread of not having macros... we can simplify code length here
	// with reflection but until then this will have to do.
	switch field {
	case "hostname":
		fmt.Printf("ip,hostname\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Hostname)
		}
	case "anycast":
		fmt.Printf("ip,anycast\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Anycast)
		}
	case "city":
		fmt.Printf("ip,city\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.City)
		}
	case "region":
		fmt.Printf("ip,region\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Region)
		}
	case "country":
		fmt.Printf("ip,country\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Country)
		}
	case "country_name":
		fmt.Printf("ip,country_name\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.CountryName)
		}
	case "loc":
		fmt.Printf("ip,loc\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Location)
		}
	case "org":
		fmt.Printf("ip,org\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Org)
		}
	case "postal":
		fmt.Printf("ip,postal\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Postal)
		}
	case "timezone":
		fmt.Printf("ip,timezone\n")
		for _, d := range core {
			fmt.Printf("%s,%v\n", d.IP, d.Timezone)
		}
	case "asn":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CoreASN{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.ASN); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "asn.id":
		fmt.Printf("ip,asn_id\n")
		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.ASN.ASN)
		}
	case "asn.name":
		fmt.Printf("ip,asn_name\n")
		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.ASN.Name)
		}
	case "asn.domain":
		fmt.Printf("ip,asn_domain\n")
		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.ASN.Domain)
		}
	case "asn.route":
		fmt.Printf("ip,asn_route\n")
		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.ASN.Route)
		}
	case "asn.type":
		fmt.Printf("ip,asn_type\n")
		for _, d := range core {
			if d.ASN == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.ASN.Type)
		}
	case "company":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CoreCompany{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Company == nil {
				continue
			}
			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.Company); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "company.name":
		fmt.Printf("ip,company_name\n")
		for _, d := range core {
			if d.Company == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Company.Name)
		}
	case "company.domain":
		fmt.Printf("ip,company_domain\n")
		for _, d := range core {
			if d.Company == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Company.Domain)
		}
	case "company.type":
		fmt.Printf("ip,company_type\n")
		for _, d := range core {
			if d.Company == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Company.Type)
		}
	case "carrier":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CoreCarrier{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.Carrier); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "carrier.name":
		fmt.Printf("ip,carrier_name\n")
		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Carrier.Name)
		}
	case "carrier.mcc":
		fmt.Printf("ip,carrier_mcc\n")
		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Carrier.MCC)
		}
	case "carrier.mnc":
		fmt.Printf("ip,carrier_mnc\n")
		for _, d := range core {
			if d.Carrier == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Carrier.MNC)
		}
	case "privacy":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CorePrivacy{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.Privacy); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "privacy.vpn":
		fmt.Printf("ip,privacy_vpn\n")
		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Privacy.VPN)
		}
	case "privacy.proxy":
		fmt.Printf("ip,privacy_proxy\n")
		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Privacy.Proxy)
		}
	case "privacy.tor":
		fmt.Printf("ip,privacy_tor\n")
		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Privacy.Tor)
		}
	case "privacy.hosting":
		fmt.Printf("ip,privacy_hosting\n")
		for _, d := range core {
			if d.Privacy == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Privacy.Hosting)
		}
	case "abuse":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CoreAbuse{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.Abuse); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "abuse.address":
		fmt.Printf("ip,abuse_address\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.Address)
		}
	case "abuse.country":
		fmt.Printf("ip,abuse_country\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,\"%v\"\n", d.IP, d.Abuse.Country)
		}
	case "abuse.country_name":
		fmt.Printf("ip,abuse_country_name\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.CountryName)
		}
	case "abuse.email":
		fmt.Printf("ip,abuse_email\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.Email)
		}
	case "abuse.name":
		fmt.Printf("ip,abuse_name\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.Name)
		}
	case "abuse.network":
		fmt.Printf("ip,abuse_network\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.Network)
		}
	case "abuse.phone":
		fmt.Printf("ip,abuse_phone\n")
		for _, d := range core {
			if d.Abuse == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Abuse.Phone)
		}
	case "domains":
		fmt.Printf("ip,")
		if err := csvEnc.EncodeHeader(ipinfo.CoreDomains{}); err != nil {
			return err
		}
		csvWriter.Flush()

		for _, d := range core {
			if d.Domains == nil {
				continue
			}

			fmt.Printf("%s,", d.IP)
			if err := csvEnc.Encode(d.Domains); err != nil {
				return err
			}
			csvWriter.Flush()
		}
	case "domains.total":
		fmt.Printf("ip,domains_total\n")
		for _, d := range core {
			if d.Domains == nil {
				continue
			}

			fmt.Printf("%s,%v\n", d.IP, d.Domains.Total)
		}
	default:
		fmt.Printf("ip,%s\n", field)
	}

	return nil
}

// Same as ipsFromCIDR with O(1) memory by discarding IPs after printing.
func outputIPsFromCIDR(cidrStr string) error {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)

	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		fmt.Println(ip)
	}

	return nil
}

// Same as ipsFromRange with O(1) memory by discarding IPs after printing.
func outputIPsFromRange(ipStrStart string, ipStrEnd string) error {
	var ipStart, ipEnd net.IP

	if ipStart = net.ParseIP(ipStrStart); ipStart == nil {
		return errNotIP
	}
	if ipEnd = net.ParseIP(ipStrEnd); ipEnd == nil {
		return errNotIP
	}

	start := binary.BigEndian.Uint32(ipStart.To4())
	end := binary.BigEndian.Uint32(ipEnd.To4())

	if start > end {
		// return decreasing list if range is flipped.
		for i := start; i >= end; i-- {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	} else {
		for i := start; i <= end; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			fmt.Println(ip)
		}
	}

	return nil
}

func createBarString(cnt int, maxCnt int) string {
	bar := "█"
	for i := 0; i < maxCnt; i++ {
		if i < cnt {
			bar += "█"
		} else {
			bar += "🮀"
		}
	}
	bar += "█"
	return bar
}

func ipsFromStdin() []net.IP {
	return ipsFromReader(os.Stdin)
}

func ipsFromReader(r io.Reader) []net.IP {
	ips := make([]net.IP, 0, 10000)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		ipStr := scanner.Text()
		if ipStr == "" {
			break
		}

		ip := net.ParseIP(ipStr)
		if ip == nil {
			// ignore any non-IP input.
			continue
		}

		ips = append(ips, ip)
	}

	return ips
}

func ipsFromCIDR(cidrStr string) ([]net.IP, error) {
	_, ipnet, err := net.ParseCIDR(cidrStr)
	if err != nil {
		return nil, err
	}

	mask := binary.BigEndian.Uint32(ipnet.Mask)
	start := binary.BigEndian.Uint32(ipnet.IP)
	end := (start & mask) | (mask ^ 0xffffffff)

	ips := make([]net.IP, 0, end-start+1)
	for i := start; i <= end; i++ {
		ip := make(net.IP, 4)
		binary.BigEndian.PutUint32(ip, i)
		ips = append(ips, ip)
	}

	return ips, nil
}

func ipsFromCIDRs(cidrStrs []string) (ips []net.IP, err error) {
	// collect IPs lists together first, then allocate a final list and do
	// a fast transfer.
	ipRanges := make([][]net.IP, len(cidrStrs))
	totalIPs := 0
	for i, cidr := range cidrStrs {
		ipRanges[i], err = ipsFromCIDR(cidr)
		if err != nil {
			return nil, err
		}
		totalIPs += len(ipRanges[i])
	}

	ips = make([]net.IP, 0, totalIPs)
	for _, ipRange := range ipRanges {
		ips = append(ips, ipRange...)
	}

	return ips, nil
}

func ipsFromRange(ipStrStart string, ipStrEnd string) ([]net.IP, error) {
	var ips []net.IP
	var ipStart, ipEnd net.IP

	if ipStart = net.ParseIP(ipStrStart); ipStart == nil {
		return nil, errNotIP
	}
	if ipEnd = net.ParseIP(ipStrEnd); ipEnd == nil {
		return nil, errNotIP
	}

	start := binary.BigEndian.Uint32(ipStart.To4())
	end := binary.BigEndian.Uint32(ipEnd.To4())
	if start > end {
		ips = make([]net.IP, 0, start-end+1)
		// return decreasing list if range is flipped.
		for i := start; i >= end; i-- {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			ips = append(ips, ip)
		}
	} else {
		ips = make([]net.IP, 0, end-start+1)
		for i := start; i <= end; i++ {
			ip := make(net.IP, 4)
			binary.BigEndian.PutUint32(ip, i)
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

func ipsFromFile(pathToFile string) ([]net.IP, error) {
	f, err := os.Open(pathToFile)
	if err != nil {
		return nil, err
	}

	return ipsFromReader(f), nil
}

func ipsFromFiles(paths []string) (ips []net.IP, err error) {
	// collect IPs lists together first, then allocate a final list and do
	// a fast transfer.
	ipLists := make([][]net.IP, len(paths))
	totalIPs := 0
	for i, p := range paths {
		ipLists[i], err = ipsFromFile(p)
		if err != nil {
			return nil, err
		}
		totalIPs += len(ipLists[i])
	}

	ips = make([]net.IP, 0, totalIPs)
	for _, ipList := range ipLists {
		ips = append(ips, ipList...)
	}

	return ips, nil
}

func saveToken(tok string) error {
	// create ipinfo config directory.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	iiCdir := filepath.Join(cdir, "ipinfo")
	if err := os.MkdirAll(iiCdir, 0700); err != nil {
		return err
	}

	// open token file.
	tokFilePath := filepath.Join(iiCdir, "token")
	tokFile, err := os.OpenFile(
		tokFilePath,
		os.O_RDWR|os.O_CREATE|os.O_TRUNC,
		0660,
	)
	defer tokFile.Close()
	if err != nil {
		return err
	}

	// write token.
	_, err = tokFile.WriteString(tok)
	if err != nil {
		return err
	}

	return nil
}

func deleteToken() error {
	// get token file path.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	tokFilePath := filepath.Join(cdir, "ipinfo", "token")

	// remove token file.
	if err := os.Remove(tokFilePath); err != nil {
		return err
	}

	return nil
}

func restoreToken() (string, error) {
	// open token file.
	cdir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	tokFilePath := filepath.Join(cdir, "ipinfo", "token")
	tokFile, err := os.Open(tokFilePath)
	defer tokFile.Close()
	if err != nil {
		return "", err
	}

	tok, err := ioutil.ReadAll(tokFile)
	if err != nil {
		return "", nil
	}

	return string(tok[:]), nil
}

func fileExists(pathToFile string) bool {
	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		return false
	}
	return true
}
