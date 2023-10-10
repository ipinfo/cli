package ipinfo

import (
	"net"
	"net/http"
	"net/netip"
)

// Core represents data from the Core API.
type Core struct {
	IP              net.IP          `json:"ip" csv:"ip"`
	Hostname        string          `json:"hostname,omitempty" csv:"hostname" yaml:"hostname,omitempty"`
	Bogon           bool            `json:"bogon,omitempty" csv:"bogon" yaml:"bogon,omitempty"`
	Anycast         bool            `json:"anycast,omitempty" csv:"anycast" yaml:"anycast,omitempty"`
	City            string          `json:"city,omitempty" csv:"city" yaml:"city,omitempty"`
	Region          string          `json:"region,omitempty" csv:"region" yaml:"region,omitempty"`
	Country         string          `json:"country,omitempty" csv:"country" yaml:"country,omitempty"`
	CountryName     string          `json:"country_name,omitempty" csv:"country_name" yaml:"countryName,omitempty"`
	CountryFlag     CountryFlag     `json:"country_flag,omitempty" csv:"country_flag_,inline" yaml:"countryFlag,omitempty"`
	CountryFlagURL  string          `json:"country_flag_url,omitempty" csv:"country_flag_url" yaml:"countryFlagURL,omitempty"`
	CountryCurrency CountryCurrency `json:"country_currency,omitempty" csv:"country_currency_,inline" yaml:"countryCurrency,omitempty"`
	Continent       Continent       `json:"continent,omitempty" csv:"continent_,inline" yaml:"continent,omitempty"`
	IsEU            bool            `json:"isEU,omitempty" csv:"isEU" yaml:"isEU,omitempty"`
	Location        string          `json:"loc,omitempty" csv:"loc" yaml:"location,omitempty"`
	Org             string          `json:"org,omitempty" csv:"org" yaml:"org,omitempty"`
	Postal          string          `json:"postal,omitempty" csv:"postal" yaml:"postal,omitempty"`
	Timezone        string          `json:"timezone,omitempty" csv:"timezone" yaml:"timezone,omitempty"`
	ASN             *CoreASN        `json:"asn,omitempty" csv:"asn_,inline" yaml:"asn,omitempty"`
	Company         *CoreCompany    `json:"company,omitempty" csv:"company_,inline" yaml:"company,omitempty"`
	Carrier         *CoreCarrier    `json:"carrier,omitempty" csv:"carrier_,inline" yaml:"carrier,omitempty"`
	Privacy         *CorePrivacy    `json:"privacy,omitempty" csv:"privacy_,inline" yaml:"privacy,omitempty"`
	Abuse           *CoreAbuse      `json:"abuse,omitempty" csv:"abuse_,inline" yaml:"abuse,omitempty"`
	Domains         *CoreDomains    `json:"domains,omitempty" csv:"domains_,inline" yaml:"domains,omitempty"`
}

// CoreASN represents ASN data for the Core API.
type CoreASN struct {
	ASN    string `json:"asn" csv:"id"`
	Name   string `json:"name" csv:"asn"`
	Domain string `json:"domain" csv:"domain"`
	Route  string `json:"route" csv:"route"`
	Type   string `json:"type" csv:"type"`
}

// CoreCompany represents company data for the Core API.
type CoreCompany struct {
	Name   string `json:"name" csv:"name"`
	Domain string `json:"domain" csv:"domain"`
	Type   string `json:"type" csv:"type"`
}

// CoreCarrier represents carrier data for the Core API.
type CoreCarrier struct {
	Name string `json:"name" csv:"name"`
	MCC  string `json:"mcc" csv:"mcc"`
	MNC  string `json:"mnc" csv:"mnc"`
}

// CorePrivacy represents privacy data for the Core API.
type CorePrivacy struct {
	VPN     bool   `json:"vpn" csv:"vpn"`
	Proxy   bool   `json:"proxy" csv:"proxy"`
	Tor     bool   `json:"tor" csv:"tor"`
	Relay   bool   `json:"relay" csv:"relay"`
	Hosting bool   `json:"hosting" csv:"hosting"`
	Service string `json:"service" csv:"service"`
}

// CoreAbuse represents abuse data for the Core API.
type CoreAbuse struct {
	Address     string `json:"address" csv:"address"`
	Country     string `json:"country" csv:"country"`
	CountryName string `json:"country_name" csv:"country_name"`
	Email       string `json:"email" csv:"email"`
	Name        string `json:"name" csv:"name"`
	Network     string `json:"network" csv:"network"`
	Phone       string `json:"phone" csv:"phone"`
}

// CoreDomains represents domains data for the Core API.
type CoreDomains struct {
	IP      string   `json:"ip" csv:"-"`
	Total   uint64   `json:"total" csv:"total"`
	Domains []string `json:"domains" csv:"-"`
}

func (v *Core) setCountryName() {
	if v.Country != "" {
		v.CountryName = GetCountryName(v.Country)
		v.IsEU = IsEU(v.Country)
		v.CountryFlag.Emoji = GetCountryFlagEmoji(v.Country)
		v.CountryFlag.Unicode = GetCountryFlagUnicode(v.Country)
		v.CountryFlagURL = GetCountryFlagURL(v.Country)
		v.CountryCurrency.Code = GetCountryCurrencyCode(v.Country)
		v.CountryCurrency.Symbol = GetCountryCurrencySymbol(v.Country)
		v.Continent.Code = GetContinentCode(v.Country)
		v.Continent.Name = GetContinentName(v.Country)
	}
	if v.Abuse != nil && v.Abuse.Country != "" {
		v.Abuse.CountryName = GetCountryName(v.Abuse.Country)
	}
}

/* CORE */

// GetIPInfo returns the details for the specified IP.
func GetIPInfo(ip net.IP) (*Core, error) {
	return DefaultClient.GetIPInfo(ip)
}

// GetIPInfoV6 returns the details for the specified IPv6 IP.
func GetIPInfoV6(ip net.IP) (*Core, error) {
	return DefaultClient.GetIPInfoV6(ip)
}

// GetIPInfo returns the details for the specified IP.
func (c *Client) GetIPInfo(ip net.IP) (*Core, error) {
	return c.getIPInfoBase(ip, false)
}

// GetIPInfoV6 returns the details for the specified IPv6 IP.
func (c *Client) GetIPInfoV6(ip net.IP) (*Core, error) {
	return c.getIPInfoBase(ip, true)
}

func (c *Client) getIPInfoBase(ip net.IP, ipv6 bool) (*Core, error) {
	relURL := ""
	if ip != nil && isBogon(netip.MustParseAddr(ip.String())) {
		bogonResponse := new(Core)
		bogonResponse.Bogon = true
		bogonResponse.IP = ip
		return bogonResponse, nil
	}
	if ip != nil {
		relURL = ip.String()
	}

	// perform cache lookup.
	if c.Cache != nil {
		if res, err := c.Cache.Get(cacheKey(relURL)); err == nil {
			return res.(*Core), nil
		}
	}

	// prepare req
	var err error
	var req *http.Request
	if ipv6 {
		req, err = c.newRequestV6(nil, "GET", relURL, nil)
	} else {
		req, err = c.newRequest(nil, "GET", relURL, nil)
	}
	if err != nil {
		return nil, err
	}

	// do req
	v := new(Core)
	if _, err := c.do(req, v); err != nil {
		return nil, err
	}

	// format
	v.setCountryName()

	// cache req result
	if c.Cache != nil {
		if err := c.Cache.Set(cacheKey(relURL), v); err != nil {
			// NOTE: still return the value even if the cache fails.
			return v, err
		}
	}

	return v, nil
}

/* IP ADDRESS */

// GetIPAddr returns the IP address that IPinfo sees when you make a request.
func GetIPAddr() (string, error) {
	return DefaultClient.GetIPAddr()
}

// GetIPAddr returns the IP address that IPinfo sees when you make a request.
func (c *Client) GetIPAddr() (string, error) {
	core, err := c.GetIPInfo(nil)
	if err != nil {
		return "", err
	}
	return core.IP.String(), nil
}

/* HOSTNAME */

// GetIPHostname returns the hostname of the domain on the specified IP.
func GetIPHostname(ip net.IP) (string, error) {
	return DefaultClient.GetIPHostname(ip)
}

// GetIPHostname returns the hostname of the domain on the specified IP.
func (c *Client) GetIPHostname(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Hostname, nil
}

/* BOGON */

// GetIPBogon returns whether an IP is a bogon IP.
func GetIPBogon(ip net.IP) (bool, error) {
	return DefaultClient.GetIPBogon(ip)
}

// GetIPBogon returns whether an IP is a bogon IP.
func (c *Client) GetIPBogon(ip net.IP) (bool, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return false, err
	}
	return core.Bogon, nil
}

/* ANYCAST */

// GetIPAnycast returns whether an IP is an anycast IP.
func GetIPAnycast(ip net.IP) (bool, error) {
	return DefaultClient.GetIPAnycast(ip)
}

// GetIPAnycast returns whether an IP is an anycast IP.
func (c *Client) GetIPAnycast(ip net.IP) (bool, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return false, err
	}
	return core.Anycast, nil
}

/* CITY */

// GetIPCity returns the city for the specified IP.
func GetIPCity(ip net.IP) (string, error) {
	return DefaultClient.GetIPCity(ip)
}

// GetIPCity returns the city for the specified IP.
func (c *Client) GetIPCity(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.City, nil
}

/* REGION */

// GetIPRegion returns the region for the specified IP.
func GetIPRegion(ip net.IP) (string, error) {
	return DefaultClient.GetIPRegion(ip)
}

// GetIPRegion returns the region for the specified IP.
func (c *Client) GetIPRegion(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Region, nil
}

/* COUNTRY */

// GetIPCountry returns the country for the specified IP.
func GetIPCountry(ip net.IP) (string, error) {
	return DefaultClient.GetIPCountry(ip)
}

// GetIPCountry returns the country for the specified IP.
func (c *Client) GetIPCountry(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Country, nil
}

/* COUNTRY NAME */

// GetIPCountryName returns the full country name for the specified IP.
func GetIPCountryName(ip net.IP) (string, error) {
	return DefaultClient.GetIPCountryName(ip)
}

// GetIPCountryName returns the full country name for the specified IP.
func (c *Client) GetIPCountryName(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.CountryName, nil
}

/* LOCATION */

// GetIPLocation returns the location for the specified IP.
func GetIPLocation(ip net.IP) (string, error) {
	return DefaultClient.GetIPLocation(ip)
}

// GetIPLocation returns the location for the specified IP.
func (c *Client) GetIPLocation(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Location, nil
}

/* ORG */

// GetIPOrg returns the organization for the specified IP.
func GetIPOrg(ip net.IP) (string, error) {
	return DefaultClient.GetIPOrg(ip)
}

// GetIPOrg returns the organization for the specified IP.
func (c *Client) GetIPOrg(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Org, nil
}

/* POSTAL */

// GetIPPostal returns the postal for the specified IP.
func GetIPPostal(ip net.IP) (string, error) {
	return DefaultClient.GetIPPostal(ip)
}

// GetIPPostal returns the postal for the specified IP.
func (c *Client) GetIPPostal(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Postal, nil
}

/* TIMEZONE */

// GetIPTimezone returns the timezone for the specified IP.
func GetIPTimezone(ip net.IP) (string, error) {
	return DefaultClient.GetIPTimezone(ip)
}

// GetIPTimezone returns the timezone for the specified IP.
func (c *Client) GetIPTimezone(ip net.IP) (string, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return "", err
	}
	return core.Timezone, nil
}

/* ASN */

// GetIPASN returns the ASN details for the specified IP.
func GetIPASN(ip net.IP) (*CoreASN, error) {
	return DefaultClient.GetIPASN(ip)
}

// GetIPASN returns the ASN details for the specified IP.
func (c *Client) GetIPASN(ip net.IP) (*CoreASN, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.ASN, nil
}

/* COMPANY */

// GetIPCompany returns the company details for the specified IP.
func GetIPCompany(ip net.IP) (*CoreCompany, error) {
	return DefaultClient.GetIPCompany(ip)
}

// GetIPCompany returns the company details for the specified IP.
func (c *Client) GetIPCompany(ip net.IP) (*CoreCompany, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.Company, nil
}

/* CARRIER */

// GetIPCarrier returns the carrier details for the specified IP.
func GetIPCarrier(ip net.IP) (*CoreCarrier, error) {
	return DefaultClient.GetIPCarrier(ip)
}

// GetIPCarrier returns the carrier details for the specified IP.
func (c *Client) GetIPCarrier(ip net.IP) (*CoreCarrier, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.Carrier, nil
}

/* PRIVACY */

// GetIPPrivacy returns the privacy details for the specified IP.
func GetIPPrivacy(ip net.IP) (*CorePrivacy, error) {
	return DefaultClient.GetIPPrivacy(ip)
}

// GetIPPrivacy returns the privacy details for the specified IP.
func (c *Client) GetIPPrivacy(ip net.IP) (*CorePrivacy, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.Privacy, nil
}

/* ABUSE */

// GetIPAbuse returns the abuse details for the specified IP.
func GetIPAbuse(ip net.IP) (*CoreAbuse, error) {
	return DefaultClient.GetIPAbuse(ip)
}

// GetIPAbuse returns the abuse details for the specified IP.
func (c *Client) GetIPAbuse(ip net.IP) (*CoreAbuse, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.Abuse, nil
}

/* DOMAINS */

// GetIPDomains returns the domains details for the specified IP.
func GetIPDomains(ip net.IP) (*CoreDomains, error) {
	return DefaultClient.GetIPDomains(ip)
}

// GetIPDomains returns the domains details for the specified IP.
func (c *Client) GetIPDomains(ip net.IP) (*CoreDomains, error) {
	core, err := c.GetIPInfo(ip)
	if err != nil {
		return nil, err
	}
	return core.Domains, nil
}
