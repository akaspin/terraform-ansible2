package inventory

import "fmt"

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
