package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/scaleoutsean/terraform-solidfire-provider/elementsw"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: elementsw.Provider,
	})
}
