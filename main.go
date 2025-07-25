package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/netapp/terraform-provider-netapp-elementsw/elementsw"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: elementsw.Provider,
	})
}
