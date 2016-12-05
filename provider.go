package terraform_ansible2

import (
	"github.com/hashicorp/terraform/terraform"
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() (p terraform.ResourceProvider) {
	
	p = &schema.Provider{
		Schema: map[string]*schema.Schema{},
		DataSourcesMap: map[string]*schema.Resource{
			"ansible2_inventory": dataInventory(),
			"ansible2_config": dataConfig(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"ansible2_play": resourcePlay(),
		},
	}
	return 
}
