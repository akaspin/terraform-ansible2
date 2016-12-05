package terraform_provider_ansible

import (
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"hash/crc32"
)

func dataPlaybook() *schema.Resource {
	return &schema.Resource{
		Read: dataPlaybookRead,
		Schema: map[string]*schema.Schema{
			"path": {
				Description: "Path to source playbook",
				Type: schema.TypeString,
				Required: true,
			},
			"include": {
				Description: "List of depended files and directories",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"exclude": {
				Description: "List of excludes",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"rendered": {
				Type: schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataPlaybookRead(d *schema.ResourceData, meta interface{}) (err error) {
	path := d.Get("path").(string)
	path, err = filepath.Abs(path)
	if err != nil {
		return
	}
	var data []byte
	data, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}

	sha := crc32.NewIEEE()
	sha.Write([]byte(path))
	sha.Write(data)

	// TODO: hash sources

	d.Set("rendered", fmt.Sprintf("# DIR %s\n# CRC32 %x\n%s", 
		filepath.Dir(path), sha.Sum(nil), string(data)))
	d.SetId(id())
	return
}

