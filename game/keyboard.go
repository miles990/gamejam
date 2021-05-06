package game

import (
	"bytes"
	"image"
	_ "image/png"
	"log"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/keyboard/keyboard"
	rkeyabord "github.com/hajimehoshi/ebiten/v2/examples/resources/images/keyboard"
)

var (
	keyboardImage *ebiten.Image
)

// Keyboard ...
type Keyboard struct {
	pressed       []ebiten.Key
	gameEventChan chan *Event
}

func (key *Keyboard) init() {
	img, _, err := image.Decode(bytes.NewReader(rkeyabord.Keyboard_png))
	if err != nil {
		log.Fatal(err)
	}

	keyboardImage = ebiten.NewImageFromImage(img)
}

// Update ...
func (key *Keyboard) Update(screen *ebiten.Image) error {
	key.pressed = nil
	for k := ebiten.Key(0); k <= ebiten.KeyMax; k++ {
		if ebiten.IsKeyPressed(k) {
			key.pressed = append(key.pressed, k)
		}
	}
	return nil
}

// Draw ...
func (key *Keyboard) Draw(screen *ebiten.Image) {
	const (
		offsetX = 24
		offsetY = 40
	)

	// Draw the base (grayed) keyboard image.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(offsetX, offsetY)
	op.ColorM.Scale(0.5, 0.5, 0.5, 1)
	screen.DrawImage(keyboardImage, op)

	// Draw the highlighted keys.
	op = &ebiten.DrawImageOptions{}
	for _, p := range key.pressed {
		op.GeoM.Reset()
		r, ok := keyboard.KeyRect(p)
		if !ok {
			continue
		}
		op.GeoM.Translate(float64(r.Min.X), float64(r.Min.Y))
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(keyboardImage.SubImage(r).(*ebiten.Image), op)
	}

	keyStrs := []string{}
	for _, p := range key.pressed {
		keyStrs = append(keyStrs, p.String())
	}
	// ebitenutil.DebugPrint(screen, strings.Join(keyStrs, ", "))
	if len(keyStrs) > 0 {
		key.sendMsg2Screen(EventKeyPress, strings.Join(keyStrs, ", "))
	}
}

func (key *Keyboard) sendMsg2Screen(eventType EventType, msg string) {
	key.gameEventChan <- &Event{
		Type: eventType,
		Msg:  msg,
	}
	log.Print(msg)
}
