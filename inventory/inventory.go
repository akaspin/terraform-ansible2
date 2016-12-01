package inventory

import (
	"fmt"
	"sort"
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
