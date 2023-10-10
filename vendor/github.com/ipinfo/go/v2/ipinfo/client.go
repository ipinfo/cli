package ipinfo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL     = "https://ipinfo.io/"
	defaultBaseURLIPv6 = "https://v6.ipinfo.io/"
	defaultUserAgent   = "IPinfoClient/Go/2.10.0"
)

// A Client is the main handler to communicate with the IPinfo API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. BaseURL should always be specified with a
	// trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the IPinfo API.
	UserAgent string

	// Cache interface implementation to prevent API quota overuse for
	// identical requests.
	Cache *Cache

	// The API token used for authorization for more data and higher limits.
	Token string
}

// NewClient returns a new IPinfo API client.
//
// If `httpClient` is nil, `http.DefaultClient` will be used.
//
// If `cache` is nil, no cache is automatically assigned. You may set one later
// at any time with `client.SetCache`.
//
// If `token` is empty, the API will be queried without any token. You may set
// one later at any time with `client.SetToken`.
func NewClient(
	httpClient *http.Client,
	cache *Cache,
	token string,
) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)
	return &Client{
		client:    httpClient,
		BaseURL:   baseURL,
		UserAgent: defaultUserAgent,
		Cache:     cache,
		Token:     token,
	}
}

// newRequest for IPV4
func (c *Client) newRequest(
	ctx context.Context,
	method string,
	urlStr string,
	body io.Reader,
) (*http.Request, error) {
	return c.newRequestBase(ctx, method, urlStr, body, false)
}

// newRequest for IPV6
func (c *Client) newRequestV6(
	ctx context.Context,
	method string,
	urlStr string,
	body io.Reader,
) (*http.Request, error) {
	return c.newRequestBase(ctx, method, urlStr, body, true)
}

// `newRequest` creates an API request. A relative URL can be provided in
// urlStr, in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
func (c *Client) newRequestBase(
	ctx context.Context,
	method string,
	urlStr string,
	body io.Reader,
	useIPv6 bool,
) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	u := new(url.URL)

	baseURL := c.BaseURL
	if useIPv6 {
		baseURL, _ = url.Parse(defaultBaseURLIPv6)
	}

	// get final URL path.
	if rel, err := url.Parse(urlStr); err == nil {
		u = baseURL.ResolveReference(rel)
	} else if strings.ContainsRune(urlStr, ':') {
		// IPv6 strings fail to parse as URLs, so let's add it as a URL Path.
		*u = *baseURL
		u.Path += urlStr
	} else {
		return nil, err
	}

	// get `http` package request object.
	req, err := http.NewRequestWithContext(ctx, method, u.String(), body)
	if err != nil {
		return nil, err
	}

	// set common headers.
	req.Header.Set("Accept", "application/json")
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	return req, nil
}

// `do` sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer interface,
// the raw response body will be written to v, without attempting to first
// decode it.
func (c *Client) do(
	req *http.Request,
	v interface{},
) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = checkResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
			if err == io.EOF {
				// ignore EOF errors caused by empty response body
				err = nil
			}
		}
	}

	return resp, err
}

// An ErrorResponse reports an error caused by an API request.
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response

	// Error structure returned by the IPinfo Core API.
	Status string `json:"status"`
	Err    struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	} `json:"error"`
}

func (r *ErrorResponse) Error() string {
	if r.Response.StatusCode == http.StatusTooManyRequests {
		return fmt.Sprintf("%v %v: %d You've hit the daily limit for the unauthenticated API. Please visit https://ipinfo.io/signup to get 50k requests per month for free.",
			r.Response.Request.Method, r.Response.Request.URL,
			r.Response.StatusCode)
	}
	return fmt.Sprintf("%v %v: %d %v",
		r.Response.Request.Method, r.Response.Request.URL,
		r.Response.StatusCode, r.Err)
}

// `checkResponse` checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range.
func checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

/* SetCache */

// SetCache assigns a cache to the package-level client.
func SetCache(cache *Cache) {
	DefaultClient.SetCache(cache)
}

// SetCache assigns a cache to the client `c`.
func (c *Client) SetCache(cache *Cache) {
	c.Cache = cache
}

/* SetToken */

// SetToken assigns a token to the package-level client.
func SetToken(token string) {
	DefaultClient.SetToken(token)
}

// SetToken assigns a token to the client `c`.
func (c *Client) SetToken(token string) {
	c.Token = token
}
