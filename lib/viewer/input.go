package viewer

import "github.com/go-gl/glfw/v3.2/glfw"

func MakeKeyCallback(v *Viewer) func(*glfw.Window, glfw.Key, int, glfw.Action, glfw.ModifierKey) {
	return func(w *glfw.Window, k glfw.Key, s int, a glfw.Action, m glfw.ModifierKey) {
		c := &called{v, w, k, s, a, m}
		v.RunIfKeyCallbacks(c)
	}
}

type called struct {
	v *Viewer
	w *glfw.Window
	k glfw.Key
	s int
	a glfw.Action
	m glfw.ModifierKey
}

type callback func(*called)

type callbacks struct {
	has map[glfw.Key][]callback
}

func newCallbacks() *callbacks {
	return &callbacks{
		make(map[glfw.Key][]callback),
	}
}

func defaultCallbacks() *callbacks {
	ret := newCallbacks()
	ret.SetKeyCallback(glfw.KeyEscape, func(c *called) {
		v := c.v
		v.Kill()
	})
	return ret
}

func (c *callbacks) GetKeyCallbacks(d *called) []callback {
	if h, ok := c.has[d.k]; ok {
		return h
	}
	return nil
}

func (c *callbacks) RunIfKeyCallbacks(d *called) {
	if h := c.GetKeyCallbacks(d); h != nil {
		for _, v := range h {
			v(d)
		}
	}
}

func (c *callbacks) SetKeyCallback(k glfw.Key, d callback) {
	if h, ok := c.has[k]; ok {
		h = append(h, d)
		c.has[k] = h
	}
	mh := make([]callback, 0)
	mh = append(mh, d)
	c.has[k] = mh
}
