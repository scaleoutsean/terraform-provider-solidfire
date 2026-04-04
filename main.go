package main

//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name solidfire

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf5to6server"
	"github.com/scaleoutsean/terraform-provider-solidfire/solidfire"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	ctx := context.Background()

	upgradedProvider, err := tf5to6server.UpgradeServer(
		ctx,
		solidfire.Provider().GRPCProvider,
	)
	if err != nil {
		log.Fatal(err)
	}

	var serveOpts []tf6server.ServeOpt
	if debug {
		serveOpts = append(serveOpts, tf6server.WithManagedDebug())
	}

	err = tf6server.Serve(
		"registry.terraform.io/scaleoutsean/solidfire",
		func() tfprotov6.ProviderServer {
			return upgradedProvider
		},
		serveOpts...,
	)

	if err != nil {
		log.Fatal(err)
	}
}
