package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rfalias/terraform-provider-pydns/pydns"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: pydns.Provider,
	})
}
