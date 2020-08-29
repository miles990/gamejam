package game

// NewGame ...
func NewGame() *Game {
	g := &Game{}
	g.init()
	return g
}

// NewTextHandler ...
func NewTextHandler() *TextHandler {
	debug := &TextHandler{}

	return debug
}

// NewHandler ...
func NewHandler(eventChan chan *Event) *EventHandler {
	handler := &EventHandler{
		gameEventChan: eventChan,
	}
	return handler
}

// NewKeyboard ...
func NewKeyboard(eventChan chan *Event) *Keyboard {
	key := &Keyboard{
		gameEventChan: eventChan,
	}
	key.init()
	return key
}
