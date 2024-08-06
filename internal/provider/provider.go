package provider

import (
	"context"
	"github.com/harvester/terraform-provider-dnsmasqapi/internal/provider/dnsentry"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/client"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider Creates the schema information for our dnsmasqapi provider and returns it to the plugin
func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			constants.FieldProviderAPIURL: {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSMASQ_API_URL", nil),
				Description: "The API Endpoint of the DNSMasq API that will service this provider.",
			},
			constants.FieldProviderAPIPort: {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("DNSMASQ_API_PORT", 0),
				Description: "The port on which the DNSMasq API listens.",
			},
			constants.FieldProviderSSLVerify: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether or not to enforce SSL validation. Set to false if using self-signed certificates.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			constants.ResourceTypeDNSEntry: dnsentry.DataSourceDnsEntry(),
		},
		ResourcesMap: map[string]*schema.Resource{
			constants.ResourceTypeDNSEntry: dnsentry.ResourceDnsEntry(),
		},
		ConfigureContextFunc: providerConfig,
	}

	return p
}

// providerConfig The configure function for the provider. Is called when the provider is configured and it's
// return value is passed as the 3rd parameter into the <Action>Context functions on Resources, so this is how we
// can inject the API client into the resource providers.
func providerConfig(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiURL := data.Get(constants.FieldProviderAPIURL).(string)
	apiPort := data.Get(constants.FieldProviderAPIPort).(int)
	sslVerify := data.Get(constants.FieldProviderSSLVerify).(bool)

	c, diags := client.NewClient(apiURL, apiPort, sslVerify)
	// TODO check DNSMasqAPI version incase these eventually drift

	return c, diags
}
