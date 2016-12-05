package terraform_provider_ansible

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() (p terraform.ResourceProvider) {
	
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"ansible_inventory": dataInventory(),
			"ansible_config": dataConfig(),
			"ansible_playbook": dataPlaybook(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ansible_playbook": resourcePlay(),
		},
	}
	return 
}
