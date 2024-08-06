package main

import (
	"github.com/harvester/terraform-provider-dnsmasqapi/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

//go:generate terraform fmt -recursive ./examples/
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.Provider}

	plugin.Serve(opts)
}
