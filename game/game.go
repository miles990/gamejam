package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Game adapter
type Game struct {
	screenWidth, screenHeight int

	drawText *TextHandler
	handler  *EventHandler
	key      *Keyboard

	isShowDebug bool

	gameEventChan chan *Event // 內部系統之間傳遞訊息用
}

var G Game

func (g *Game) init() *Game {
	g.screenWidth = 320
	g.screenHeight = 240

	g.isShowDebug = true

	g.drawText = NewTextHandler()
	g.drawText.init(g.isShowDebug)

	g.gameEventChan = make(chan *Event, 1000)

	g.handler = NewHandler(g.gameEventChan)
	go g.handler.Run()

	g.key = NewKeyboard(g.gameEventChan)

	return g
}

// Update : Update the logical state
func (g *Game) Update() error {
	// if err := g.key.Update(screen); err != nil {
	// 	return err
	// }
	g.key.pressed = inpututil.PressedKeys()
	return nil
}

// Draw : Render the screen
func (g *Game) Draw(screen *ebiten.Image) {

	var event *Event
	select {
	case event = <-g.gameEventChan:
	default:
	}

	g.drawText.Draw(screen, event, g.screenWidth, g.screenHeight)
	g.key.Draw(screen)

	// fmt.Println(time.Now())

}

// Layout : Return the game screen size
func (g *Game) Layout(outsideWidth, outsigeHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}
