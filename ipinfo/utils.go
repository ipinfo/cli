package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/ipinfo/go/v2/ipinfo"
)

func isIP(ipStr string) bool {
	return net.ParseIP(ipStr) != nil
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
		fmt.Printf("Country         %s\n", d.Abuse.Country)
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
