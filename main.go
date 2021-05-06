package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/miles990/gamejam/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("title")
	// g := game.NewGame()

	// game.G.Init()
	if err := ebiten.RunGame(&game.G); err != nil {
		panic(err)
	}
}
