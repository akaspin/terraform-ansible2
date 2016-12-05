package terraform_ansible2

import (
	"crypto/rand"
	"fmt"
	"crypto/sha1"
	"encoding/hex"
)

func id() (r string) {
	b := make([]byte, 6)
	rand.Read(b)
	r = fmt.Sprintf("%x", b)
	return 
}

func hashId(what interface{}) (r string) {
	hash := sha1.Sum([]byte(what.(string)))
	r = hex.EncodeToString(hash[:])
	return 
}
