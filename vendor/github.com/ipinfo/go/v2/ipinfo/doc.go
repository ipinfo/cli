/*
Package ipinfo provides a client for using the IPinfo API.

Usage:

	import "github.com/ipinfo/go/v2/ipinfo"

The default IPinfo client is predefined and can be used without initialization.
For example:

	info, err := ipinfo.GetIPInfo(net.ParseIP("8.8.8.8"))

Authorization

To perform authorized API calls with more data and higher limits, pass in a
non-empty token to NewClient. For example:

	client := ipinfo.NewClient(nil, nil, "MY_TOKEN")
	info, err := client.GetIPInfo(net.ParseIP("8.8.8.8"))
*/
package ipinfo
