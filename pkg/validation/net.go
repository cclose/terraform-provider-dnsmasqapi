package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/harvester/terraform-provider-dnsmasqapi/pkg/constants"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

var Validator *validator.Validate

func init() {
	Validator = validator.New(validator.WithRequiredStructEnabled())
}

// ValidateHostname checks if the provided hostname is valid
func ValidateHostname(val interface{}, path cty.Path) diag.Diagnostics {
	hostname, ok := val.(string)
	if !ok {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid hostname type",
			Detail:   "Hostname must be a string",
		}}
	}

	err := Validator.Var(hostname, constants.ValidateHostnameRequired)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// ValidateIPAddress checks if the provided IP Address is valid
func ValidateIPAddress(val interface{}, path cty.Path) diag.Diagnostics {
	ip, ok := val.(string)
	if !ok {
		return diag.Diagnostics{diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Invalid IP type",
			Detail:   "IP must be a string",
		}}
	}

	err := Validator.Var(ip, constants.ValidateIPRequired)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil

}
