package viewer

import (
	"fmt"
	"image"
	"os"
	"runtime"

	"github.com/Laughs-In-Flowers/log"
	i "github.com/Laughs-In-Flowers/valhalla/lib/image"
	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

type Viewer struct {
	*glfw.Window
	log.Logger
	*callbacks
	initialized bool
	kill        bool
	image       *i.Image
}

func New(path string, l log.Logger) (*Viewer, error) {
	a := defaultCallbacks()
	im, err := i.New(path)
	if err != nil {
		return &Viewer{nil, l, a, false, false, nil}, err
	}
	l.Printf("opened image at %s", path)
	v := &Viewer{nil, l, a, false, false, im}
	err = initialize(v)
	if err != nil {
		return nil, err
	}
	return v, nil
}

func bound(im image.Image) (int, int) {
	var w, h int
	const maxSize = 1200
	w = im.Bounds().Size().X
	h = im.Bounds().Size().Y
	a := float64(w) / float64(h)
	if a >= 1 {
		if w > maxSize {
			w = maxSize
			h = int(maxSize / a)
		}
	} else {
		if h > maxSize {
			h = maxSize
			w = int(maxSize * a)
		}
	}
	return w, h
}

func initialize(v *Viewer) error {
	if !v.initialized {
		if err := gl.Init(); err != nil {
			return err
		}

		if err := glfw.Init(); err != nil {
			return err
		}

		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 5)
		glfw.WindowHint(glfw.Resizable, glfw.False)

		w, h := bound(v.image)

		ip := v.image.Path()
		wt := fmt.Sprintf("valhalla -- %s", ip)

		window, err := glfw.CreateWindow(w, h, wt, nil, nil)
		if err != nil {
			return err
		}

		window.MakeContextCurrent()
		glfw.SwapInterval(1)

		window.SetKeyCallback(MakeKeyCallback(v))

		v.Window = window
		v.SetRefreshCallback(v.onRefresh)

		gl.ClearColor(0.11, 0.545, 0.765, 0.0)

		v.initialized = true
	}

	return nil
}

func (v *Viewer) Run() {
	v.Print("viewing...")
	runtime.LockOSThread()
RUN:
	for {
		if v.kill {
			break RUN
		}
		v.Draw()
		glfw.PollEvents()
	}
	close(v)
}

func close(v *Viewer) {
	v.image.Free()
	v.Destroy()
	glfw.Terminate()
	v.Print("exiting")
	os.Exit(0)
}

func (v *Viewer) Kill() {
	v.kill = true
}

func (v *Viewer) onRefresh(*glfw.Window) {
	v.Draw()
}

func (v *Viewer) Draw() {
	v.MakeContextCurrent()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	v.image.Draw()
	v.SwapBuffers()
}
