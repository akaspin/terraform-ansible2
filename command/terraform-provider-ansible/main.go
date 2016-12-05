package main

import (
	"github.com/akaspin/terraform-provider-ansible"
	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
	"log"
	"github.com/hashicorp/terraform/helper/logging"
)

var V string

func main() {
	out, _ := logging.LogOutput()
	log.SetOutput(out)
	
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return terraform_provider_ansible.Provider()
		},
	})
}
