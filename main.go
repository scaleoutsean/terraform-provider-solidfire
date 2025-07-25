package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scaleoutsean/terraform-provider-solidfire/elementsw"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: elementsw.Provider,
	})
}
