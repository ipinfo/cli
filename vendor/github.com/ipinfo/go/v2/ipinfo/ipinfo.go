package ipinfo

// DefaultClient is the package-level client available to the user.
var DefaultClient *Client

func init() {
	/* Create a global, default client. */
	DefaultClient = NewClient(nil, nil, "")
}
