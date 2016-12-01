package terraform_ansible2

import (
	"github.com/hashicorp/terraform/helper/schema"
	"io/ioutil"
	"fmt"
	"path/filepath"
	"os"
	"hash/crc32"
)

func dataPlaybook() *schema.Resource {
	return &schema.Resource{
		Read: dataPlaybookRead,
		Schema: map[string]*schema.Schema{
			"contents": {
				Description: "Playbook contents",
				Type: schema.TypeString,
				Optional: true,
			},
			"path": {
				Description: "Path to playbook",
				Type: schema.TypeString,
				Required: true,
			},
			"sources": {
				Description: "List of files and directories in playbook",
				Type: schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"rendered": {
				Type: schema.TypeString,
				Computed: true,
			},
			"directory": {
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
	stat, err := os.Stat(path)
	if err != nil {
		return 
	}
	if stat.IsDir() {
		d.Set("directory", path)
	} else {
		d.Set("directory", filepath.Dir(path))
	}
	
	// if contents is absent use path to fetch playbook
	var contents string
	raw, ok := d.GetOk("contents")
	if !ok {
		if stat.IsDir() {
			err = fmt.Errorf("%s is directory", path)
			return 
		}
		var data []byte
		data, err = ioutil.ReadFile(path)
		if err != nil {
			return 
		}
		contents = string(data)
	} else {
		contents = raw.(string)
	}
	
	sha := crc32.NewIEEE()
	sha.Write([]byte(path))
	sha.Write([]byte(contents))
	
	// TODO: hash sources
	
	h := fmt.Sprintf("%x", sha.Sum(nil))
	d.Set("rendered", fmt.Sprintf("# CRC32 %s\n%s", h, contents))
	d.SetId(h)
	return 
}
