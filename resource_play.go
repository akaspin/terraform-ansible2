package terraform_ansible2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/akaspin/terraform-ansible2/ansible"
	"github.com/hashicorp/terraform/helper/logging"
	"os"
)

// resourcePlay  
func resourcePlay() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"playbook": {
				Description: "Playbook path",
				Type: schema.TypeString,
				Required: true,
			},
			"inventory": {
				Description: "Inventory contents",
				Type: schema.TypeString,
				Required: true,
			},
			"config": {
				Description: "Config contents",
				Type: schema.TypeString,
				Optional: true,
			},
			"extra_json": {
				Description: "Extra variables JSON",
				Type: schema.TypeString,
				Optional: true,
			},
			"limit": {
				Description: "Limit play",
				Type: schema.TypeString,
				Optional: true,
				Default: "all",
			},
			"tags": {
				Description: "Tags",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"skip_tags": {
				Description: "Skip tags",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"phase": {
				Description: "Phase control",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create": {
							Description: "Run on create",
							Type:     schema.TypeBool,
							Optional: true,
							Default: true,
						},
						"update": {
							Description: "Run on modify",
							Type:     schema.TypeBool,
							Optional: true,
							Default: true,
						},
						"destroy": {
							Description: "Run on destroy",
							Type:     schema.TypeBool,
							Optional: true,
							Default: false,
						},
						"tag": {
							Description: "Add create/modify/destroy tag",
							Type: schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"untagged": {
							Description: "Add untagged phase tag",
							Type: schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
				MaxItems: 1,
			},
			"cleanup": {
				Description: "Remove ansible files after successful run",
				Type: schema.TypeBool,
				Optional: true,
				Default: false,
			},
		},
		Create: resourcePlayCreate,
		Update: resourcePlayUpdate,
		Read: resourcePlayRead,
		Delete: resourcePlayDelete,
	}
}

func resourcePlayCreate(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlayRead(d, meta)
	d.SetId(id())
	if runner, ok := resourcePlayGetRunner(d, meta, "create"); ok {
		err = runner.Run()
	}
	return 
}

func resourcePlayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlayRead(d, meta)
	//d.Partial(true)
	
	if runner, ok := resourcePlayGetRunner(d, meta, "update"); ok {
		err = runner.Run()
	} 
	return 
}

func resourcePlayRead(d *schema.ResourceData, meta interface{}) (err error) {
	return 
}

func resourcePlayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlayRead(d, meta)
	runner, ok := resourcePlayGetRunner(d, meta, "destroy")
	if ok {
		if err = runner.Run(); err != nil {
			return 
		}
	}
	runner.Cleanup()
	d.SetId("")
	return 
}

func resourcePlayGetRunner(d *schema.ResourceData, meta interface{}, phase string) (r *ansible.Playbook, ok bool) {
	output, err := logging.LogOutput()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	phaseOpts := resourcePlayPhase(d)
	
	tags := extractStringSlice(d, "tags")
	if phaseOpts["tag"] {
		tags = append(tags, phase)
	}
	if phaseOpts["untagged"] {
		tags = append(tags, "untagged")
	}
	
	config := ansible.PlaybookConfig{
		Id: d.Id(),
		Config: d.Get("config").(string),
		ExtraJson: d.Get("extra_json").(string),
		Inventory: d.Get("inventory").(string),
		PlaybookPath: d.Get("playbook").(string),
		Tags: tags,
		SkipTags: extractStringSlice(d, "skip_tags"),
		Limit: d.Get("limit").(string),
		CleanupOnSuccess: d.Get("cleanup").(bool),
	}
	r = ansible.NewPlaybook(wd, output, config)
	ok = phaseOpts[phase]
	return 
}

func resourcePlayPhase(d *schema.ResourceData) (r map[string]bool) {
	r = map[string]bool{}
	if raw, ok := d.GetOk("phase"); ok {
		for k, v := range raw.([]interface{})[0].(map[string]interface{}) {
			r[k] = v.(bool)
		}
		return 
	}
	r = map[string]bool{
		"create": true,
		"update": true,
		"destroy": false,
		"tag": true,
		"untagged": true,
	}
	return 
}
