package inventory

import "fmt"

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
