package dnsentry

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSourceDnsEntry() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDnsEntryRead,
		Schema:      DataSourceSchema(),
	}
}
