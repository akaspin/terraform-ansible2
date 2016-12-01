package terraform_ansible2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"fmt"
	"github.com/akaspin/terraform-ansible2/ansible"
	"io"
	"github.com/hashicorp/terraform/helper/logging"
	"hash/crc32"
)

// resourcePlay  
func resourcePlay() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"playbook": {
				Description: "Playbook contents",
				Type: schema.TypeString,
				Required: true,
			},
			"directory": {
				Description: "Playbook directory",
				Type: schema.TypeString,
				Optional: true,
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
				Description: "Extra variables",
				Type: schema.TypeString,
				Optional: true,
				Default: "'{}'",
			},
			"limit": {
				Description: "Limit play",
				Type: schema.TypeString,
				Optional: true,
				Default: "all",
			},
			"on_destroy": {
				Description: "Run on destroy",
				Type: schema.TypeBool,
				Optional: true,
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
	if err = resourcePlayRead(d, meta); err != nil {
		d.SetId("")
		return 
	}
	
	runner, err := resourcePlayGetRunner(d, meta, "create")
	if err != nil {
		return 
	}
	
	if err = runner.Run(); err != nil {
		d.SetId("")
		return 
	}
	return 
}

func resourcePlayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	// cleanup
	output, err := getOutput()
	if err != nil {
		return 
	}
	prev_directory, _ := d.GetChange("directory")
	runner, err := ansible.NewPlaybook(output, ansible.PlaybookConfig{
		Id: d.Id(),
		PlayDir: prev_directory.(string),
	})
	if err != nil {
		return 
	}
	runner.Cleanup()
	
	if err = resourcePlayRead(d, meta); err != nil {
		return 
	}
	runner, err = resourcePlayGetRunner(d, meta, "update")
	if err != nil {
		return
	}

	if err = runner.Run(); err != nil {
		d.SetId("")
		return 
	}
	return 
}

func resourcePlayRead(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlaySetId(d)
	return 
}

func resourcePlayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	runner, err := resourcePlayGetRunner(d, meta, "destroy")
	if err != nil {
		return 
	}
	if d.Get("on_destroy").(bool) {
		if err = runner.Run(); err != nil {
			d.SetId("")
			return 
		}
	}
	runner.Cleanup()
	d.SetId("")
	return 
}

func resourcePlayGetRunner(d *schema.ResourceData, meta interface{}, phase string) (r *ansible.Playbook, err error) {
	config := ansible.PlaybookConfig{
		Id: d.Id(),
		Config: d.Get("config").(string),
		Extra: d.Get("extra_json").(string),
		Inventory: d.Get("inventory").(string),
		Playbook: d.Get("playbook").(string),
		PlayDir: d.Get("directory").(string),
		Limit: d.Get("limit").(string),
		Phase: phase,
		CleanupOnSuccess: d.Get("cleanup").(bool),
	}
	output, err := getOutput()
	if err != nil {
		return 
	}
	r, err = ansible.NewPlaybook(output, config)
	return 
}

func resourcePlaySetId(d *schema.ResourceData) (err error) {
	sha := crc32.NewIEEE()
	sha.Write([]byte(d.Get("playbook").(string)))
	sha.Write([]byte(d.Get("directory").(string)))
	sha.Write([]byte(d.Get("inventory").(string)))
	sha.Write([]byte(d.Get("config").(string)))
	sha.Write([]byte(d.Get("extra_json").(string)))
	sha.Write([]byte(d.Get("limit").(string)))
	sha.Write([]byte(fmt.Sprintf("%t", d.Get("on_destroy").(bool))))
	d.SetId(fmt.Sprintf("%x", sha.Sum(nil)))
	return 
}

func dump(what string, d *schema.ResourceData) {
	msg := fmt.Sprintf("%s[%t] %s", what, d.IsNewResource(), d.Id())
	for _, k := range []string{
		"playbook",
		"inventory",
		"config",
		"extra_json",
		"limit",
		"root",
	} {
		msg += fmt.Sprintf(" %s(%t)", k, d.HasChange(k))
	}
	log.Printf("%s", msg)
	log.Printf(">>> %s", d.Get("root"))
}

func getOutput() (io.Writer, error) {
	return logging.LogOutput()
}
