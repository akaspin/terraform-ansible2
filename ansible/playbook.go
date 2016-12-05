package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PlaybookConfig struct {
	Id               string
	PlaybookPath     string
	Inventory        string
	Config           string
	ExtraJson        string
	Tags []string
	SkipTags []string
	Limit            string
	Phase            string
	CleanupOnSuccess bool
}

type Playbook struct {
	WD     string
	Output io.Writer
	PlaybookConfig
}

func NewPlaybook(output io.Writer, config PlaybookConfig) (p *Playbook, err error) {
	p = &Playbook{
		Output:         output,
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
	}
	if extra != "" {
		params = append(params, []string{"-e", fmt.Sprintf("'%s'", extra)}...)
	}
	if len(p.Tags) > 0 {
		params = append(params, []string{
			"--tags", 
			fmt.Sprintf(`"%s"`, strings.Join(p.Tags, ",")),
		}...)
	}
	if len(p.SkipTags) > 0 {
		params = append(params, []string{
			"--skip-tags", 
			fmt.Sprintf(`"%s"`, strings.Join(p.SkipTags, ",")),
		}...)
	}

	cmd := exec.Command("ansible-playbook", append(params, p.PlaybookPath)...)
	cmd.Env = append(os.Environ(),
		[]string{
			fmt.Sprintf("ANSIBLE_LOG_PATH=%s", p.assetPath("log")),
			"ANSIBLE_RETRY_FILES_ENABLED=no",
		}...)
	if p.Config != "" {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("ANSIBLE_CONFIG=%s", p.assetPath("cfg")))
	}

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
	os.Remove(p.assetPath("cfg"))
	os.Remove(p.assetPath("inventory"))
	os.Remove(p.assetPath("log"))
}

func (p *Playbook) prepare() (err error) {
	p.Cleanup()
	if err = ioutil.WriteFile(p.assetPath("inventory"), []byte(p.Inventory), 0755); err != nil {
		return
	}
	if p.Config != "" {
		if err = ioutil.WriteFile(p.assetPath("cfg"), []byte(p.Config), 0755); err != nil {
			return
		}
	}
	return
}

func (p *Playbook) extra() (r string, err error) {
	if p.ExtraJson == "" {
		return 
	}
	mid := map[string]interface{}{}
	if err = json.Unmarshal([]byte(p.ExtraJson), &mid); err != nil {
		return
	}

	data, err := json.Marshal(&mid)
	if err != nil {
		return
	}
	r = string(data)
	return
}

func (p *Playbook) assetPath(kind string) (r string) {
	r = filepath.Join(p.WD, fmt.Sprintf(".ansible-%s.%s", p.Id, kind))
	return
}
