package terraform_provider_ansible

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/logging"
	"os"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"strings"
)

func resourcePlay() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"playbook": {
				Description: "Prepared playbook",
				Type: schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"playbook_path"},
			},
			"playbook_path": {
				Description: "Playbook path",
				Type: schema.TypeString,
				Optional: true,
				ConflictsWith: []string{"playbook"},
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
	d.SetId(id())
	if runner, ok := resourcePlayGetRunner(d, meta, "create"); ok {
		err = runner.Run()
	}
	return 
}

func resourcePlayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	if runner, ok := resourcePlayGetRunner(d, meta, "update"); ok {
		err = runner.Run()
	} 
	return 
}

func resourcePlayRead(d *schema.ResourceData, meta interface{}) (err error) {
	return 
}

func resourcePlayDelete(d *schema.ResourceData, meta interface{}) (err error) {
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

func resourcePlayGetRunner(d *schema.ResourceData, meta interface{}, phase string) (r *Play, ok bool) {
	output, err := logging.LogOutput()
	if err != nil {
		panic(err)
	}
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	
	// prepare playbook
	var playbookContents, playbookDir string
	
	if rawPath, rawExists := d.GetOk("playbook_path"); rawExists {
		var data []byte
		if data, err = ioutil.ReadFile(rawPath.(string)); err != nil {
			return 
		}
		playbookContents = string(data)
		playbookDir = filepath.Dir(rawPath.(string))
	} else if rawContents, rawExists1 := d.GetOk("playbook"); rawExists1 {
		playbookContents = rawContents.(string)
		for _, line := range strings.Split(playbookContents, "\n") {
			if strings.HasPrefix(line, "# DIR ") {
				playbookDir = strings.TrimPrefix(line, "# DIR ")
				break
			}
		}
		if playbookDir == "" {
			err = fmt.Errorf(
				`No directory in prepared playbook. Use "ansible_playbook" datasource.`)
		}
		
	} else {
		err = fmt.Errorf("One of playbook or playbook_path is required")
		return 
	}
	
	
	phaseOpts := map[string]bool{}
	if raw, exists := d.GetOk("phase"); exists {
		for k, v := range raw.([]interface{})[0].(map[string]interface{}) {
			phaseOpts[k] = v.(bool)
		}
	} else {
		phaseOpts = map[string]bool{
			"create": true,
			"update": true,
			"destroy": false,
			"tag": true,
			"untagged": true,
		}
	}
	
	tags := extractStringSlice(d, "tags")
	if phaseOpts["tag"] {
		tags = append(tags, phase)
	}
	if phaseOpts["untagged"] {
		tags = append(tags, "untagged")
	}
	
	r = &Play{
		WD: wd,
		Output: output,
		Id: d.Id(),
		Config: d.Get("config").(string),
		ExtraJson: d.Get("extra_json").(string),
		Inventory: d.Get("inventory").(string),
		Playbook: playbookContents,
		PlayDir: playbookDir,
		Tags: tags,
		SkipTags: extractStringSlice(d, "skip_tags"),
		Limit: d.Get("limit").(string),
		CleanupOnSuccess: d.Get("cleanup").(bool),
	}
	
	ok = phaseOpts[phase]
	return 
}
