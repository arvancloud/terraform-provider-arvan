package main

import (
	"flag"
	"github.com/arvancloud/terraform-provider-arvan/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"log"
)

func init() {
	logFlags := log.Flags()
	logFlags = logFlags &^ (log.Ldate | log.Ltime)
	log.SetFlags(logFlags)
}

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
		Debug:        debugMode,
	})
}
