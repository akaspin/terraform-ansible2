package terraform_provider_ansible

import (
	"github.com/hashicorp/terraform/helper/schema"
	"fmt"
	"sort"
)

func dataConfig() *schema.Resource {
	return &schema.Resource{
		Read: dataConfigRead,
		Schema: dataConfigSchema(),
	}
}

func dataConfigHierarchy() map[string]map[string]*schema.Schema {
	return map[string]map[string]*schema.Schema{
		"defaults": map[string]*schema.Schema{
			"roles_path": {
				Type: schema.TypeString,
				Optional: true,
			},
			"library": {
				Type: schema.TypeString,
				Optional: true,
			},
			"filter_plugins": {
				Type: schema.TypeString,
				Optional: true,
			},
			"forks": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"display_skipped_hosts": {
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
			"stdout_callback": {
				Type: schema.TypeString,
				Optional: true,
			},
			"callback_whitelist": {
				Type: schema.TypeString,
				Optional: true,
			},
			
			"transport": {
				Type: schema.TypeString,
				Optional: true,
			},
			"timeout": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"remote_port": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"poll_interval": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"host_key_checking": {
				Type: schema.TypeBool,
				Optional: true,
			},
			
			"gathering": {
				Type: schema.TypeString,
				Optional: true,
			},
			"gather_subset": {
				Type: schema.TypeString,
				Optional: true,
			},
			"gather_timeout": {
				Type: schema.TypeInt,
				Optional: true,
			},
			
			"remote_user": {
				Type: schema.TypeString,
				Optional: true,
			},
			"private_key_file": {
				Type: schema.TypeString,
				Optional: true,
			},
			"sudo_flags": {
				Type: schema.TypeString,
				Optional: true,
			},
			
			"jinja2_extensions": {
				Type: schema.TypeString,
				Optional: true,
			},
			
			"task_includes_static": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"handler_includes_static": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"module_lang": {
				Type: schema.TypeInt,
				Optional: true,
			},
			"module_set_locale": {
				Type: schema.TypeBool,
				Optional: true,
			},
		},
		"privilege_escalation": map[string]*schema.Schema{
			"become": {
				Type: schema.TypeBool,
				Optional: true,
			},
			"become_method": {
				Type: schema.TypeString,
				Optional: true,
			},
			"become_user": {
				Type: schema.TypeString,
				Optional: true,
			},
		},
		"ssh_connection": map[string]*schema.Schema{
			"ssh_args": {
				Type: schema.TypeString,
				Optional: true,
			},
			"control_path_dir": {
				Type: schema.TypeString,
				Optional: true,
			},
			"control_path": {
				Type: schema.TypeString,
				Optional: true,
			},
			"pipelining": {
				Type: schema.TypeBool,
				Optional: true,
			},
		},
		"computed": map[string]*schema.Schema{
			"rendered": {
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataConfigSchema() (r map[string]*schema.Schema) {
	r = map[string]*schema.Schema{}
	for _, section := range dataConfigHierarchy() {
		for k, v := range section {
			r[k] = v
		}
	} 
	return 
}

func dataConfigSortedKeys(section string) (r []string) {
	var t sort.StringSlice
	for k := range dataConfigHierarchy()[section] {
		t = append(t, k)
	}
	t.Sort()
	r = t
	return 
}

func dataConfigRead(d *schema.ResourceData, meta interface{}) (err error) {
	sections := []string{
		"defaults", "privilege_escalation", "ssh_connection",
	}
	var r string
	for _, section := range sections {
		r += fmt.Sprintf("[%s]\n", section)
		for _, key := range dataConfigSortedKeys(section) {
			if v, ok := d.GetOk(key); ok {
				r += fmt.Sprintf("%s = %v\n", key, v)
			}
		}
		r += "\n"
	} 
	
	d.Set("rendered", r)
	d.SetId(hash(r))
	return 
}
