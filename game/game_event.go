package game

const ()

// EventType ...
type EventType int

// EventType enum
const (
	EventConnect EventType = iota
	EventError
	EventMessage
	EventDisconnect
	EventServerSubscribe
	EventServerUnsubscribe
	EventServerJoin
	EventServerLeave
	EventServerPublish
	EventPublish
	EventJoin
	EventLeave
	EventSubscribeSuccess
	EventSubscribeError
	EventUnsubscribe

	EventKeyPress

	EventClearInfoText
)

var (
	eventName = map[EventType]string{
		EventConnect:           "connect",
		EventError:             "error",
		EventMessage:           "msg",
		EventDisconnect:        "disconnect",
		EventServerSubscribe:   "server_subscribe",
		EventServerUnsubscribe: "server_unsubscribe",
		EventServerJoin:        "server_join",
		EventServerLeave:       "server_leave",
		EventServerPublish:     "server_publish",
		EventPublish:           "publish",
		EventJoin:              "join",
		EventLeave:             "leave",
		EventSubscribeSuccess:  "subscribe_success",
		EventSubscribeError:    "subscribe_error",
		EventUnsubscribe:       "unsubscribe",
		EventKeyPress:          "key_press",
		EventClearInfoText:     "clear_info_text",
	}
)

func (e EventType) String() string {
	return eventName[e]
}

// Event ...
type Event struct {
	Type EventType
	Msg  string
}
