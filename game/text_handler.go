package game

import (
	"bytes"
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// DrawFuncType ...
type DrawFuncType int

// DrawFuncType enum
const (
	DrawFPS DrawFuncType = iota
	DrawCusor
	DrawInfo
)

// // DebugInfo num
// const (
// 	DebugInfoNum = 3
// )

// EventMsgFunc ...
type EventMsgFunc func(e ...*Event) string

// TextHandler ...
type TextHandler struct {
	drawTextFuncMap map[DrawFuncType]EventMsgFunc
	lastInfoText    string
}

func (d *TextHandler) init(isDebug bool) {
	d.drawTextFuncMap = make(map[DrawFuncType]EventMsgFunc, 0)

	d.drawTextFuncMap[DrawFPS] = func(e ...*Event) string {
		fps := ebiten.CurrentFPS()
		return fmt.Sprintf("fps: %.2f", fps)
	}

	d.drawTextFuncMap[DrawCusor] = func(e ...*Event) string {
		cursorX, cursorY := ebiten.CursorPosition()
		return fmt.Sprintf("cusorX: %d, cursorY:%d", cursorX, cursorY)
	}

	d.drawTextFuncMap[DrawInfo] = func(e ...*Event) string {
		if len(e) == 0 {
			return d.lastInfoText
		}
		var event = e[0]
		if event == nil {
			return d.lastInfoText
		}
		if event.Type == EventClearInfoText {
			return ""
		}

		var buffer bytes.Buffer
		n := 200
		var n1 = n - 1
		var l1 = len(event.Msg) - 1
		for i, rune := range event.Msg {
			buffer.WriteRune(rune)
			if i%n == n1 && i != l1 {
				buffer.WriteRune('\n')
			}
		}
		d.lastInfoText = fmt.Sprintf("[%s]\n%s", event.Type, buffer.String())
		return d.lastInfoText
	}
}

// Draw ...
func (d *TextHandler) Draw(screen *ebiten.Image, e *Event, screenW int, screenH int) {
	fps := d.drawTextFuncMap[DrawFPS]()
	cursor := d.drawTextFuncMap[DrawCusor]()

	ebitenutil.DebugPrint(screen, fmt.Sprintf("%s\n%s", fps, cursor))
	info := d.drawTextFuncMap[DrawInfo](e)
	ebitenutil.DebugPrintAt(screen, info, 0, 200)
}

// // Draw ...
// func (d *DebugHandler) Draw(screen *ebiten.Image) {
// 	// d.gameEventChan <- &Event{
// 	// 	Type: eventType,
// 	// 	Msg:  msg,
// 	// }
// 	// log.Print(msg)
// }
