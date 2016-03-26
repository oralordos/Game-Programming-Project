package main

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

type Input struct {
	X, Y float64
}

func (i *Input) Normalize() {
	length := math.Sqrt(i.X*i.X + i.Y*i.Y)
	i.X /= length
	i.Y /= length
}

func (i *Input) Combine(other *Input) {
	i.X += other.X
	i.Y += other.Y
}

type InputSystem interface {
	ProcessEvent(sdl.Event, *PlayerFrontend) (bool, *Input)
}

type ExitInput struct{}

func (e ExitInput) ProcessEvent(ev sdl.Event, front *PlayerFrontend) (bool, *Input) {
	if _, ok := ev.(*sdl.QuitEvent); ok {
		front.Destroy()
		return true, nil
	}
	return false, nil
}

type KeyboardInput struct {
	left, right, up, down bool
	a, d, w, s            bool
}

func (k *KeyboardInput) ProcessEvent(ev sdl.Event, front *PlayerFrontend) (bool, *Input) {
	switch e := ev.(type) {
	case *sdl.KeyDownEvent:
		if e.Repeat == 0 {
			switch e.Keysym.Sym {
			case sdl.K_a:
				k.a = true
			case sdl.K_LEFT:
				k.left = true
			case sdl.K_d:
				k.d = true
			case sdl.K_RIGHT:
				k.right = true
			case sdl.K_w:
				k.w = true
			case sdl.K_UP:
				k.up = true
			case sdl.K_s:
				k.s = true
			case sdl.K_DOWN:
				k.down = true
			}
		}
	case *sdl.KeyUpEvent:
		switch e.Keysym.Sym {
		case sdl.K_a:
			k.a = false
		case sdl.K_LEFT:
			k.left = false
		case sdl.K_d:
			k.d = false
		case sdl.K_RIGHT:
			k.right = false
		case sdl.K_w:
			k.w = false
		case sdl.K_UP:
			k.up = false
		case sdl.K_s:
			k.s = false
		case sdl.K_DOWN:
			k.down = false
		}
	}
	return false, k.createStruct()
}

func (k *KeyboardInput) createStruct() *Input {
	i := Input{}
	if k.left || k.a {
		i.X--
	}
	if k.right || k.d {
		i.X++
	}
	if k.up || k.w {
		i.Y--
	}
	if k.down || k.s {
		i.Y++
	}
	return &i
}
