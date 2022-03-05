package main

import (
	"embed"
	"github.com/go-gl/glfw/v3.3/glfw"
)

//go:embed assets
var assets embed.FS

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	game := &Game{}
	window := &Window{InitialSize: Dimensions{Width: 900, Height: 520}}
	game.Window = window
	window.Create(game)

	for !window.ShouldClose() {
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
