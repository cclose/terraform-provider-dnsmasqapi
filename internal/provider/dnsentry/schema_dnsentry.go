package dnsentry

import (
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/constants"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/validation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		constants.FieldDNSEntryHostname: {
			Type:             schema.TypeString,
			Required:         true,
			Description:      constants.FieldDescDNSEntryHostname,
			ValidateDiagFunc: validation.ValidateHostname,
		},
		constants.FieldDNSEntryIpAddress: {
			Type:             schema.TypeString,
			Required:         true,
			Description:      constants.FieldDescDNSEntryIpAddress,
			ValidateDiagFunc: validation.ValidateIPAddress,
		},
		constants.FieldDNSEntryAppend: {
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
			Description: `Add this IP Address to this hostname instead of overriding the current value.`,
		},
	}
}

func DataSourceSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		constants.FieldDNSEntryHostname: {
			Type:             schema.TypeString,
			Required:         true,
			Description:      constants.FieldDescDNSEntryHostname,
			ValidateDiagFunc: validation.ValidateHostname,
		},
		constants.FieldDNSEntryIpAddress: {
			Type:        schema.TypeString,
			Computed:    true,
			Description: constants.FieldDescDNSEntryIpAddress,
		},
	}
}
