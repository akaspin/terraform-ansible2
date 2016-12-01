package ansible

import (
	"os"
	"path/filepath"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"io"
	"os/exec"
	"log"
)

type PlaybookConfig struct {
	Id        string
	Inventory string
	Playbook  string
	PlayDir   string
	Config    string
	Extra     string
	Limit     string
	Phase     string
	CleanupOnSuccess   bool
} 

type Playbook struct {
	WD string
	Output io.Writer
	PlaybookConfig	
}

func NewPlaybook(output io.Writer, config PlaybookConfig) (p *Playbook, err error) {
	p = &Playbook{
		Output: output,
		PlaybookConfig: config,
	}
	p.WD, err = os.Getwd()
	return 
}

func (p *Playbook) Run() (err error) {
	var extra string
	if extra, err = p.extra(); err != nil {
		return 
	}
	if err = p.prepare(); err != nil {
		return 
	}
	
	params := []string{
		"-l", p.Limit,
		"-i", p.assetPath("inventory"),
		"-e", fmt.Sprintf("'%s'", extra),
		p.playbookPath(),
	}
	
	cmd := exec.Command("ansible-playbook", params...)
	cmd.Env = append(os.Environ(), 
		[]string{
			fmt.Sprintf("ANSIBLE_CONFIG=%s", p.assetPath("cfg")),
			fmt.Sprintf("ANSIBLE_LOG_PATH=%s", p.assetPath("log")),
			"ANSIBLE_RETRY_FILES_ENABLED=no",
		}...)
	log.Printf("running ansible-playbook %s %s", params, cmd.Env)
	cmd.Stdout = p.Output
	cmd.Stderr = p.Output
	if err = cmd.Run(); err != nil {
		err = fmt.Errorf("%v : see log %s", err, p.assetPath("log"))
		return 
	}
	if p.CleanupOnSuccess {
		p.Cleanup()
	}
	return 
}

// Silently cleanup all files
func (p *Playbook) Cleanup() {
	os.Remove(p.playbookPath())
	os.Remove(p.assetPath("cfg"))
	os.Remove(p.assetPath("inventory"))
	os.Remove(p.assetPath("log"))
}

func (p *Playbook) prepare() (err error) {
	p.Cleanup()
	if err = ioutil.WriteFile(p.playbookPath(), []byte(p.Playbook), 0755); err != nil {
		return 
	}
	if err = ioutil.WriteFile(p.assetPath("cfg"), []byte(p.Config), 0755); err != nil {
		return 
	}
	if err = ioutil.WriteFile(p.assetPath("inventory"), []byte(p.Inventory), 0755); err != nil {
		return 
	}
	if err = ioutil.WriteFile(p.assetPath("log"), []byte(p.Inventory), 0755); err != nil {
		return 
	}
	return 
}

func (p *Playbook) extra() (r string, err error) {
	mid := map[string]interface{}{}
	if err = json.Unmarshal([]byte(p.Extra), &mid); err != nil {
		return 
	}
	mid["phase"] = p.Phase
	
	data, err := json.Marshal(&mid)
	if err != nil {
		return 
	}
	r = string(data)
	return 
}

func (p *Playbook) assetPath(kind string) (r string) {
	r = filepath.Join(p.WD, fmt.Sprintf(".tf-%s.%s", p.Id, kind))
	return 
}

func (p *Playbook) playbookPath() (r string) {
	r = filepath.Join(p.PlayDir, fmt.Sprintf(".tf-%s.yaml", p.Id))
	return 
}
