package constants

// Validation
const (
	ValidateRequired         = "required"
	ValidateHostname         = "hostname_rfc1123"
	ValidateIP               = "ip"
	ValidateHostnameRequired = ValidateRequired + "," + ValidateHostname
	ValidateIPRequired       = ValidateRequired + "," + ValidateIP
)

// Provider
const (
	FieldProviderAPIURL    = "api_url"
	FieldProviderAPIPort   = "api_port"
	FieldProviderSSLVerify = "ssl_verify"
)

// Resource DNSEntry
const (
	ResourceTypeDNSEntry = "dnsmasqapi_dnsentry"

	FieldDNSEntryAppend    = "append"
	FieldDNSEntryHostname  = "hostname"
	FieldDNSEntryIpAddress = "ip_address"

	FieldDescDNSEntryHostname  = "The hostname of the DNS entry. Should be a Fully Qualified Domain Name (FQDN)."
	FieldDescDNSEntryIpAddress = "The IP Addess of the DNS entry."
)
