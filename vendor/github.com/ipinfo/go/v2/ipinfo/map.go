package ipinfo

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
)

// IPMap is the full JSON response from the IP Map API.
type IPMap struct {
	Status    string `json:"status"`
	ReportURL string `json:"reportUrl"`
}

// GetIPMap returns an IPMap result for a group of IPs.
//
// `len(ips)` must not exceed 500,000.
func GetIPMap(ips []net.IP) (*IPMap, error) {
	return DefaultClient.GetIPMap(ips)
}

// GetIPMap returns an IPMap result for a group of IPs.
//
// `len(ips)` must not exceed 500,000.
func (c *Client) GetIPMap(ips []net.IP) (*IPMap, error) {
	if len(ips) > 500000 {
		return nil, errors.New("ip count must be <500,000")
	}

	jsonArrStr, err := json.Marshal(ips)
	if err != nil {
		return nil, err
	}
	jsonBuf := bytes.NewBuffer(jsonArrStr)

	req, err := c.newRequest(nil, "POST", "map?cli=1", jsonBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	result := new(IPMap)
	if _, err := c.do(req, result); err != nil {
		return nil, err
	}

	return result, nil
}
