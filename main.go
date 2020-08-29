package main

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/miles990/gamejam/game"
)

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("title")
	g := game.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
