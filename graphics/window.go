package graphics

import (
	"runtime"

	"github.com/veandco/go-sdl2/sdl"
)

func Init() error {
	// This is needed to make sure that we are always running the graphics code in the same OS thread as it is initialized in.
	runtime.LockOSThread()
	if runtime.GOMAXPROCS(-1) == 1 {
		runtime.GOMAXPROCS(2)
	}
	return sdl.Init(sdl.INIT_EVERYTHING)
}

func Quit() {
	sdl.Quit()
	runtime.UnlockOSThread()
}

type Window struct {
	win  *sdl.Window
	rend *sdl.Renderer
}

func CreateWindow(width, height int, title string) (*Window, error) {
	var g Window
	var err error
	g.win, err = sdl.CreateWindow(title, sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, width, height, 0)
	if err != nil {
		return nil, err
	}
	g.rend, err = sdl.CreateRenderer(g.win, -1, 0)
	return &g, err
}

func (g *Window) Destroy() {
	g.rend.Destroy()
	g.win.Destroy()
}

func (g *Window) Clear() error {
	g.rend.SetDrawColor(0, 0, 0, 255)
	return g.rend.Clear()
}

func (g *Window) Draw(d Drawable) error {
	return d.Draw(g.rend, 0, 0)
}

func (g *Window) Update(draw Drawable) error {
	if err := g.Clear(); err != nil {
		return err
	}
	if err := g.Draw(draw); err != nil {
		return err
	}
	g.Present()
	return nil
}

func (g *Window) Present() {
	g.rend.Present()
}

func (g *Window) GetSize() (int, int) {
	return g.win.GetSize()
}
