package terraform_provider_ansible

import (
	"github.com/hashicorp/terraform/helper/schema"
	"hash/crc32"
	"fmt"
	"strings"
)

func dataInventory() *schema.Resource {
	return &schema.Resource{
		Read: dataInventoryRead,
		Schema: map[string]*schema.Schema{
			"hosts": {
				Description: "Inventory hosts",
				Type: schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Description: "Hosts group",
							Required: true,
							Type: schema.TypeString,
						},
						"names": {
							Description: "Inventory host names in order",
							Required: true,
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							MinItems: 1,
						},
						"vars": {
							Description: "Group variables. May prepend with `cast:type`",
							Optional: true,
							Type: schema.TypeMap,
							Elem: &schema.Schema{
								Type: schema.TypeString,	
							},
						},
					},
				},
			},
			"var": {
				Description: "Host bounded variables in order",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Description: "Host group",
							Type: schema.TypeString,
							Required: true,
						},
						"key": {
							Description: "Variable key",
							Type: schema.TypeString,
							Required: true,
						},
						"values": {
							Description: "Values in order",
							Type: schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cast": {
							Description: "Cast value (string, int, float, bool)",
							Type: schema.TypeString,
							Optional: true,
							Default: "string",
						},
					},
				},
			},

			"rendered": {
				Type: schema.TypeString,
				Computed: true,
			},

		},
	}
}

func dataInventoryRead(d *schema.ResourceData, meta interface{}) (err error) {
	i := NewInventory()
	
	for _, raw := range (d.Get("hosts")).([]interface{}) {
		chunk := raw.(map[string]interface{})
		group := (chunk["group"]).(string)
		i.AddGroup(group)

		var hostnames []string
		for _, hostname := range chunk["names"].([]interface{}) {
			hostnames = append(hostnames, hostname.(string))
		}
		if err = i.AddHosts(group, hostnames...); err != nil {
			return
		}

		for key, v := range (chunk["vars"]).(map[string]interface {}) {
			value := strings.TrimSpace(v.(string))
			cast := "string"
			if strings.HasPrefix(value, "`cast:") {
				value = strings.TrimPrefix(value, "`cast:")
				split := strings.SplitN(value, "`", 2)
				cast = strings.TrimSpace(split[0])
				value = strings.TrimSpace(split[1])
			}
			var v1 *Variable
			if v1, err = NewVariable(key, value, cast); err != nil {
				return
			}
			i.AddGroupVar(group, v1)
		}
	}

	for _, raw := range (d.Get("var")).([]interface{}) {
		chunk := raw.(map[string]interface{})
		group := (chunk["group"]).(string)
		name := (chunk["key"]).(string)
		cast := (chunk["cast"]).(string)

		var variables []*Variable
		for _, val := range chunk["values"].([]interface{}) {
			var v *Variable
			if v, err = NewVariable(name, val.(string), cast); err != nil {
				return
			}
			variables = append(variables, v)
		}
		if err = i.BindHostVars(group, variables...); err != nil {
			return
		}
	}

	r, err := i.Render()
	if err != nil {
		return
	}
	d.Set("rendered", r)
	d.SetId(hash(r))

	return
}

func hash(s string) string {
	sha := crc32.NewIEEE()
	sha.Write([]byte(s))
	return fmt.Sprintf("%x", sha.Sum(nil))
}
