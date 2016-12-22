package terraform_provider_ansible

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

type Play struct {
	WD     string
	Output io.Writer

	Id               string
	Playbook         string
	Inventory        string
	Config           string
	ExtraJson        string
	Tags             []string
	SkipTags         []string
	Limit            string
	Verbosity string
	CleanupOnSuccess bool
}

func (p *Play) Run() (err error) {
	var extra string
	if extra, err = p.extra(); err != nil {
		return
	}
	
	
	params := []string{
		"-l", p.Limit,
		"-i", p.assetPath("inventory"),
	}
	if p.Verbosity != "" {
		params = append(params, fmt.Sprintf("-%s", p.Verbosity))
	}
	
	if extra != "" {
		params = append(params, []string{"-e", fmt.Sprintf("%s", extra)}...)
	}
	if len(p.Tags) > 0 {
		params = append(params, []string{
			"--tags", 
			fmt.Sprintf(`%s`, strings.Join(p.Tags, ",")),
		}...)
	}
	if len(p.SkipTags) > 0 {
		params = append(params, []string{
			"--skip-tags", 
			fmt.Sprintf(`%s`, strings.Join(p.SkipTags, ",")),
		}...)
	}
	
	cmd := exec.Command("ansible-playbook", append(params, p.Playbook)...)
	cmd.Env = append(os.Environ(),
		[]string{
			fmt.Sprintf("ANSIBLE_LOG_PATH=%s", p.assetPath("log")),
			"ANSIBLE_RETRY_FILES_ENABLED=no",
			"ANSIBLE_NOCOLOR=1",
		}...)
	if p.Config != "" {
		cmd.Env = append(cmd.Env,
			fmt.Sprintf("ANSIBLE_CONFIG=%s", p.assetPath("cfg")))
	}
	
	// add run options to log
	loglines := strings.Join([]string{
		//strings.Join(cmd.Env, " "),
		"\nansible-playbook",
		strings.Join(cmd.Args, " "),
		"\n",
	}, "")
	
	p.Cleanup()
	if err = p.writeAsset("inventory", p.Inventory); err != nil {
		return
	}
	if p.Config != "" {
		if err = p.writeAsset("cfg", p.Config); err != nil {
			return
		}
	}
	if err = p.writeAsset("log", loglines); err != nil {
		return
	}

	fmt.Fprintf(p.Output, "%#v %t", p.Output, p.Output)
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
func (p *Play) Cleanup() {
	os.Remove(p.assetPath("cfg"))
	os.Remove(p.assetPath("inventory"))
	os.Remove(p.assetPath("log"))
}

func (p *Play) extra() (r string, err error) {
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

func (p *Play) assetPath(kind string) (r string) {
	r = filepath.Join(".ansible", fmt.Sprintf("%s.%s", p.Id, kind))
	return
}

func (p *Play) writeAsset(kind string, data string) (err error) {
	ap := p.assetPath(kind)
	if err = os.MkdirAll(".ansible", 0755); err != nil {
		return 
	}
	err = ioutil.WriteFile(ap, []byte(data), 0755)
	return 
}
