package ipinfo

import (
	"strings"
)

// ASNDetails represents details for an ASN.
type ASNDetails struct {
	ASN         string             `json:"asn"`
	Name        string             `json:"name"`
	Country     string             `json:"country"`
	CountryName string             `json:"-"`
	Allocated   string             `json:"allocated"`
	Registry    string             `json:"registry"`
	Domain      string             `json:"domain"`
	NumIPs      uint64             `json:"num_ips"`
	Type        string             `json:"type"`
	Prefixes    []ASNDetailsPrefix `json:"prefixes"`
	Prefixes6   []ASNDetailsPrefix `json:"prefixes6"`
	Peers       []string           `json:"peers"`
	Upstreams   []string           `json:"upstreams"`
	Downstreams []string           `json:"downstreams"`
}

// ASNDetailsPrefix represents data for prefixes managed by an ASN.
type ASNDetailsPrefix struct {
	Netblock string `json:"netblock"`
	ID       string `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Size     string `json:"size"`
	Status   string `json:"status"`
	Domain   string `json:"domain"`
}

// InvalidASNError is reported when the invalid ASN was specified.
type InvalidASNError struct {
	ASN string
}

func (err *InvalidASNError) Error() string {
	return "invalid ASN: " + err.ASN
}

func (v *ASNDetails) setCountryName() {
	if v.Country != "" {
		v.CountryName = countriesMap[v.Country]
	}
}

// GetASNDetails returns the details for the specified ASN.
func GetASNDetails(asn string) (*ASNDetails, error) {
	return DefaultClient.GetASNDetails(asn)
}

// GetASNDetails returns the details for the specified ASN.
func (c *Client) GetASNDetails(asn string) (*ASNDetails, error) {
	if !strings.HasPrefix(asn, "AS") {
		return nil, &InvalidASNError{ASN: asn}
	}

	// perform cache lookup.
	if c.Cache != nil {
		if res, err := c.Cache.Get(cacheKey(asn)); err == nil {
			return res.(*ASNDetails), nil
		}
	}

	// prepare req
	req, err := c.newRequest(nil, "GET", asn, nil)
	if err != nil {
		return nil, err
	}

	// do req
	v := new(ASNDetails)
	if _, err := c.do(req, v); err != nil {
		return nil, err
	}

	// format
	v.setCountryName()

	// cache req result
	if c.Cache != nil {
		if err := c.Cache.Set(cacheKey(asn), v); err != nil {
			return v, err
		}
	}

	return v, nil
}
