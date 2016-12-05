package terraform_ansible2

import (
	"crypto/rand"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
	"github.com/hashicorp/terraform/helper/schema"
)

func id() (r string) {
	b := make([]byte, 6)
	rand.Read(b)
	r = fmt.Sprintf("%x", b)
	return 
}

func extractStringSlice(d *schema.ResourceData, k string) (r []string) {
	raw, ok := d.GetOk(k)
	if !ok {
		return 
	}
	for _, v := range raw.([]interface{}) {
		r = append(r, v.(string))
	}
	return 
}

func hashId(what interface{}) (r string) {
	hash := sha1.Sum([]byte(what.(string)))
	r = hex.EncodeToString(hash[:])
	return 
}
