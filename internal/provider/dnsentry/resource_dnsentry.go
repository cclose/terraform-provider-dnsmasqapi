package dnsentry

import (
	"context"
	"fmt"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/client"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/constants"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDnsEntry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsEntryCreate,
		ReadContext:   resourceDnsEntryRead,
		UpdateContext: resourceDnsEntryUpdate,
		DeleteContext: resourceDnsEntryDelete,
		Schema:        ResourceSchema(),
	}
}

func resourceDnsEntryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	hostname := d.Get(constants.FieldDNSEntryHostname).(string)
	ipAddress := d.Get(constants.FieldDNSEntryIpAddress).(string)
	apnd := d.Get(constants.FieldDNSEntryAppend).(bool)

	_, err := apiClient.PostDNSEntry(hostname, ipAddress, apnd)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO verify returned entries match request?

	return resourceDnsEntryRead(ctx, d, m)
}

func resourceDnsEntryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	hostname := d.Get(constants.FieldDNSEntryHostname).(string)
	// Implement your API call to read the DNS entry
	// ...
	entries, err := apiClient.GetDNSEntry(hostname)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO, we need to change ip into a list, i think...
	d.SetId(fmt.Sprintf("%s|%s", hostname, entries[0].IP)) // Use hostname|ip as the ID
	d.Set(constants.FieldDescDNSEntryIpAddress, entries[0].IP)

	return nil
}

func resourceDnsEntryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	hostname := d.Get(constants.FieldDNSEntryHostname).(string)
	ipAddress := d.Get(constants.FieldDNSEntryIpAddress).(string)
	apnd := d.Get(constants.FieldDNSEntryAppend).(bool)

	_, err := apiClient.PostDNSEntry(hostname, ipAddress, apnd)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO verify returned entries match request?

	return resourceDnsEntryRead(ctx, d, m)
}

func resourceDnsEntryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	apiClient := m.(*client.Client)
	hostname := d.Get(constants.FieldDNSEntryHostname).(string)

	err := apiClient.DeleteDNSEntry(hostname)
	if err != nil {
		diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
