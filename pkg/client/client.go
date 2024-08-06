package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cclose/dnsmasq-api/model"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	httpClient *http.Client
	APIAddress string
	SSLVerify  bool
}

// NewClient creates a new DNSMasqAPI Client for resources to use
func NewClient(apiURL string, apiPort int, sslVerify bool) (*Client, diag.Diagnostics) {
	var diags diag.Diagnostics
	// Check if the URL has a scheme (http or https)
	if !strings.HasPrefix(apiURL, "http://") && !strings.HasPrefix(apiURL, "https://") {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "URL Lacks Protocol",
			Detail:   fmt.Sprintln("API URL lacks 'http://' or 'https://'. Defaulting to 'https://'."),
		})
		apiURL = "https://" + apiURL
	}

	// Parse the URL
	parsedURL, err := url.Parse(apiURL)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("invalid API URL: %w", err))
	}

	// Check if the URL already contains a port
	if parsedURL.Port() != "" {
		return nil, diag.FromErr(fmt.Errorf("API URL should not contain a port. Use the separate field " +
			constants.FieldProviderAPIPort,
		))
	}

	// Construct the API address with the port if provided
	apiAddress := apiURL
	if apiPort != 0 {
		apiAddress = fmt.Sprintf("%s:%d", apiURL, apiPort)
	}

	// Create and return the client
	client := &Client{
		APIAddress: apiAddress,
		SSLVerify:  sslVerify,
	}

	return client, diags
}

// doRequest sends an HTTP request to the API
func (c *Client) doRequest(method string, route string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.APIAddress, route)

	var reqBody []byte
	var err error
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Configure SSL verification
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.SSLVerify},
	}
	client := &http.Client{Transport: transport}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d and body: %v", res.StatusCode, res.Body)
	}

	return res, nil
}

// decodeDNSRecords decodes a JSON response into DNSRecords
func (c *Client) decodeDNSRecords(res *http.Response) ([]model.DNSRecord, error) {
	var entries []model.DNSRecord
	if res.Body != nil {
		if res.Header.Get("Content-Type") == "application/json" {
			err := json.NewDecoder(res.Body).Decode(&entries)
			if err != nil {
				return nil, fmt.Errorf("failed to parse response body: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unexpected response content-type: %s", res.Header.Get("Content-Type"))
		}
	}

	return entries, nil
}

// DeleteDNSEntry deletes all DNS entries for the specified hostname
func (c *Client) DeleteDNSEntry(hostname string) error {
	if hostname == "" {
		return fmt.Errorf("hostname must be provided")
	}

	route := fmt.Sprintf("dns/%s", hostname)
	// delete returns a string so we don't care, status code tells us all we need
	_, err := c.doRequest(http.MethodDelete, route, nil)

	return err
}

// GetAllDNSEntries retrieves all DNS entries on the API
func (c *Client) GetAllDNSEntries() ([]model.DNSRecord, error) {
	route := "dns"
	res, err := c.doRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	return c.decodeDNSRecords(res)
}

// GetDNSEntry retrieves all entries for the specified hostname
func (c *Client) GetDNSEntry(hostname string) ([]model.DNSRecord, error) {
	if hostname == "" {
		return nil, fmt.Errorf("hostname must be provided")
	}

	route := fmt.Sprintf("dns/%s", hostname)
	res, err := c.doRequest(http.MethodGet, route, nil)
	if err != nil {
		return nil, err
	}

	return c.decodeDNSRecords(res)
}

// PostDNSEntry sends a POST request to the API to update or create a DNS entry. Will replace existing
// entries unless ?append is set
func (c *Client) PostDNSEntry(hostname string, ip string, append bool) ([]model.DNSRecord, error) {
	if hostname == "" || ip == "" {
		return nil, fmt.Errorf("hostname and ip must be provided")
	}

	route := fmt.Sprintf("dns/%s", hostname)
	if append {
		route += "?append=true"
	}

	// map ipaddress into the proper format
	body := model.SetDNSRecordRequest{
		IPs: []string{ip},
	}

	res, err := c.doRequest(http.MethodPost, route, body)
	if err != nil {
		return nil, err
	}

	return c.decodeDNSRecords(res)
}
