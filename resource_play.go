package terraform_ansible2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/akaspin/terraform-ansible2/ansible"
	"github.com/hashicorp/terraform/helper/logging"
	"io"
	"log"
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
				StateFunc: hashId,
			},
			"config": {
				Description: "Config contents",
				Type: schema.TypeString,
				Optional: true,
				StateFunc: hashId,
			},
			"extra_json": {
				Description: "Extra variables JSON",
				Type: schema.TypeString,
				Optional: true,
				StateFunc: hashId,
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
			//"phase": {
			//	Description: "Phase control",
			//	Type: schema.TypeList,
			//	Optional: true,
			//	//Computed: true,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"create": {
			//				Description: "Run on create",
			//				Type:     schema.TypeBool,
			//				Optional: true,
			//				Default:  true,
			//			},
			//			"modify": {
			//				Description: "Run on modify",
			//				Type:     schema.TypeBool,
			//				Optional: true,
			//				Default:  true,
			//			},
			//			"destroy": {
			//				Description: "Run on destroy",
			//				Type:     schema.TypeBool,
			//				Optional: true,
			//				Default:  false,
			//			},
			//			
			//			"tag": {
			//				Description: "Add create/modify/destroy tag",
			//				Type: schema.TypeBool,
			//				Optional: true,
			//				Default: true,
			//			},
			//			"strict_tag": {
			//				Description: "Do not add untagged to tags if phase tag is on",
			//				Type: schema.TypeBool,
			//				Optional: true,
			//				Default: false,
			//			},
			//		},
			//	},
			//	MaxItems: 1,
			//},
			"cleanup": {
				Description: "Remove ansible files after successful run",
				Type: schema.TypeBool,
				Optional: true,
				//Default: false,
			},
			"test": {
				Description: "Remove ansible files after successful run",
				Type: schema.TypeInt,
				Optional: true,
				//Default: false,
			},
			"test_string": {
				Description: "Remove ansible files after successful run",
				Type: schema.TypeString,
				Optional: true,
				//Default: false,
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
	runner, err := resourcePlayGetRunner(d, meta, "create")
	if err != nil {
		d.SetId("")
		return 
	}
	
	if err = runner.Run(); err != nil {
		//d.SetId("")
		return 
	}
	return 
}

func resourcePlayUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlayRead(d, meta)
	d.Partial(true)
	
	runner, err := resourcePlayGetRunner(d, meta, "update")
	if err != nil {
		return
	}
	runner.Cleanup()

	d.SetId(id())
	runner, err = resourcePlayGetRunner(d, meta, "update")
	if err != nil {
		return
	}
	err = runner.Run()
	return 
}

func resourcePlayRead(d *schema.ResourceData, meta interface{}) (err error) {
	log.Printf(">>> read %#v", d.Get("phase"))
	//if _, ok := d.GetOk("phase"); !ok {
	//	if err := d.Set("phase", []interface{}{}); err != nil {
	//		return err
	//	}
	//}
	
	
	return 
}

func resourcePlayDelete(d *schema.ResourceData, meta interface{}) (err error) {
	resourcePlayRead(d, meta)
	runner, err := resourcePlayGetRunner(d, meta, "destroy")
	if err != nil {
		return 
	}
	//if err = runner.Run(); err != nil {
	//	return 
	//}
	runner.Cleanup()
	d.SetId("")
	return 
}

func resourcePlayGetRunner(d *schema.ResourceData, meta interface{}, phase string) (r *ansible.Playbook, err error) {
	log.Print("create runner")
	
	
	config := ansible.PlaybookConfig{
		Id: d.Id(),
		Config: d.Get("config").(string),
		ExtraJson: d.Get("extra_json").(string),
		Inventory: d.Get("inventory").(string),
		PlaybookPath: d.Get("playbook").(string),
		Limit: d.Get("limit").(string),
		Phase: phase,
		CleanupOnSuccess: d.Get("cleanup").(bool),
	}
	output, err := getOutput()
	if err != nil {
		return 
	}
	r, err = ansible.NewPlaybook(output, config)
	if err == nil {
		log.Print("runner created")
	}
	return 
}

func getOutput() (io.Writer, error) {
	return logging.LogOutput()
}
