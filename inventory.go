package terraform_provider_ansible

import (
	"fmt"
	"sort"
	"encoding/json"
)

type Inventory struct {
	Hosts map[string]*Host
	Groups map[string]*Group
}

func NewInventory() (i *Inventory) {
	i = &Inventory{
		Hosts: map[string]*Host{},
		Groups: map[string]*Group{},
	}
	return
}

func (i *Inventory) AddGroup(name string) {
	if _, exists := i.Groups[name]; exists {
		return
	}
	i.Groups[name] = NewGroup(name)
}


func (i *Inventory) AddHosts(group string, hosts ...string) (err error) {
	var g *Group
	if g, err = i.getGroup(group); err != nil {
		err = fmt.Errorf("can't add hosts : %s", err)
		return
	}

	for _, hostname := range hosts {
		var h *Host
		var exists bool
		if h, exists = i.Hosts[hostname]; !exists {
			i.Hosts[hostname] = NewHost(hostname)
			h = i.Hosts[hostname]
		}
		g.AddHost(h)
	}
	return
}

func (i *Inventory) BindHostVars(group string, variables ...*Variable) (err error) {
	var g *Group
	if g, err = i.getGroup(group); err != nil {
		err = fmt.Errorf("can't bind host variables : %s", err)
	}
	err = g.BindHostVariables(variables...)
	return
}

func (i *Inventory) AddGroupVar(group string, variable *Variable) (err error) {
	var g *Group
	if g, err = i.getGroup(group); err != nil {
		err = fmt.Errorf("can't add group variable : %s", err)
	}
	g.AddVariable(variable)
	return
}

func (i *Inventory) getGroup(name string) (g *Group, err error) {
	var ok bool
	if g, ok = i.Groups[name]; !ok {
		err = fmt.Errorf("group %s is not found", name)
	}
	return
}

func (i *Inventory) Render() (r string, err error) {
	r = "[all]\n"
	
	var hostNames sort.StringSlice
	for n := range i.Hosts {
		hostNames = append(hostNames, n)
	}
	hostNames.Sort()

	for _, host := range hostNames {
		var chunk string
		if chunk, err = i.Hosts[host].Render(); err != nil {
			return
		}
		r += fmt.Sprintf("%s\n", chunk)
	}
	r += "\n"

	var groupNames sort.StringSlice
	for n := range i.Groups {
		groupNames = append(groupNames, n)
	}
	groupNames.Sort()

	for _, group := range groupNames {
		var chunk string
		if chunk, err = i.Groups[group].Render(); err != nil {
			return
		}
		r += fmt.Sprintf("%s\n", chunk)
	}
	return
}


type Group struct {
	Name string
	Hosts []*Host
	Variables []*Variable
}

func NewGroup(name string) (g *Group) {
	g = &Group{
		Name: name,
	}
	return
}

func (g *Group) AddHost(host *Host) {
	for _, h := range g.Hosts {
		if h.Name == host.Name {
			return
		}
	}
	g.Hosts = append(g.Hosts, host)
	return
}

func (g *Group) BindHostVariables(variables ...*Variable) (err error) {
	if len(variables) != len(g.Hosts) {
		err = fmt.Errorf("bad vars list length hosts=%d vars =%d", len(g.Hosts), len(variables))
		return
	}
	for i, host := range g.Hosts {
		host.AddVariable(variables[i])
	}
	return
}

func (g *Group) AddVariable(variable *Variable) {
	g.Variables = append(g.Variables, variable)
}

func (g *Group) Render() (r string, err error) {
	r = fmt.Sprintf("[%s]\n", g.Name)
	for _, host := range g.Hosts {
		r += fmt.Sprintf("%s\n", host.Name)
	}
	if len(g.Variables) > 0 {
		r += fmt.Sprintf("\n[%s:vars]\n", g.Name)
		for _, v := range g.Variables {
			var v1 string
			if v1, err = v.Render(); err != nil {
				return
			}
			r += fmt.Sprintf("%s\n", v1)
		}
	}
	return
}

type Host struct {
	Name string
	Variables []*Variable
}

func NewHost(name string) (h *Host) {
	h = &Host{
		Name: name,
	}
	return
}

func (h *Host) AddVariable(v *Variable) {
	h.Variables = append(h.Variables, v)
}

func (h *Host) Render() (r string, err error) {
	r = h.Name
	for _, v := range h.Variables {
		var rendered string
		if rendered, err = v.Render(); err != nil {
			return
		}
		r += fmt.Sprintf(" %s", rendered)
	}
	return
}

// Variable represents zero-level inventory variable
type Variable struct {

	// Variable name
	Name string

	// Variable raw value
	Value string

	// Variable output cast
	CastFunc CastFunc
}

func NewVariable(name, value, cast string) (v *Variable, err error) {
	v = &Variable{
		Name: name,
		Value: value,
	}
	switch cast {
	case "string":
		v.CastFunc = CastString
	case "json":
		v.CastFunc = CastJson
	}
	return
}

func (v *Variable) Render() (r string, err error) {
	var value string
	if value, err = v.CastFunc(v.Value); err != nil {
		return
	}
	r = fmt.Sprintf("%s=%s", v.Name, value)
	return
}

type CastFunc func(in string) (r string, err error)

func CastString(in string) (r string, err error) {
	r = fmt.Sprintf("%s", in)
	return
}

func CastJson(in string) (r string, err error) {
	var tmp interface{}
	if err = json.Unmarshal([]byte(in), &tmp); err != nil {
		return
	}
	var data []byte
	if data, err = json.Marshal(&tmp); err != nil {
		return
	}
	r = fmt.Sprintf("'%s'", string(data))
	return
}
