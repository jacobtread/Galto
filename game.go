package main

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type Game struct {
	Window *Window
}

func (g *Game) onMouseMove(window *glfw.Window, x float64, y float64) {

}

func (g *Game) onMouseClick(window *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {

}

func (g Game) onKeyPressed(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

}

func (g *Game) onCharPressed(window *glfw.Window, char rune) {

}
