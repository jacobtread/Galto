package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
)

type Dimensions struct {
	Width  int
	Height int
}

type Window struct {
	*glfw.Window

	VSync       bool
	InitialSize Dimensions
}

func init() {
	runtime.LockOSThread()
}

func (w *Window) Create(game *Game) {

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DepthBits, 24)
	window, err := glfw.CreateWindow(w.InitialSize.Width, w.InitialSize.Height, "Galto", nil, nil)
	if err != nil {
		panic(err)
	}
	w.Window = window
	window.MakeContextCurrent()

	window.SetCursorPosCallback(game.onMouseMove)
	window.SetMouseButtonCallback(game.onMouseClick)
	window.SetKeyCallback(game.onKeyPressed)
	window.SetCharCallback(game.onCharPressed)
}
