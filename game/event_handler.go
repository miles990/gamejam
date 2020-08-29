package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	// _ "net/http/pprof"
	"github.com/centrifugal/centrifuge-go"
	"github.com/dgrijalva/jwt-go"
)

// In real life clients should never know secret key. This is only for example
// purposes to quickly generate JWT for connection.
const exampleTokenHmacSecret = "8a9d4508-a038-4b66-8a44-5820cb305334"

func connToken(user string, exp int64) string {
	// NOTE that JWT must be generated on backend side of your application!
	// Here we are generating it on client side only for example simplicity.
	claims := jwt.MapClaims{"sub": user}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(exampleTokenHmacSecret))
	if err != nil {
		panic(err)
	}
	return t
}

// ChatMessage is chat example specific message struct.
type ChatMessage struct {
	Input string `json:"input"`
}

// EventHandler ...
type EventHandler struct {
	gameEventChan chan *Event
}

// OnConnect ...
func (h *EventHandler) OnConnect(_ *centrifuge.Client, e centrifuge.ConnectEvent) {
	h.sendMsg2Screen(EventConnect, fmt.Sprintf("Connected ID %s", e.ClientID))
}

// OnError ...
func (h *EventHandler) OnError(_ *centrifuge.Client, e centrifuge.ErrorEvent) {
	h.sendMsg2Screen(EventError, fmt.Sprintf("Error: %s", e.Message))
}

// OnMessage ...
func (h *EventHandler) OnMessage(_ *centrifuge.Client, e centrifuge.MessageEvent) {
	h.sendMsg2Screen(EventMessage, fmt.Sprintf("Message from server: %s", string(e.Data)))
}

// OnDisconnect ...
func (h *EventHandler) OnDisconnect(_ *centrifuge.Client, e centrifuge.DisconnectEvent) {
	h.sendMsg2Screen(EventDisconnect, fmt.Sprintf("Disconnected from chat: %s", e.Reason))
}

// OnServerSubscribe ...
func (h *EventHandler) OnServerSubscribe(_ *centrifuge.Client, e centrifuge.ServerSubscribeEvent) {
	h.sendMsg2Screen(EventServerSubscribe, fmt.Sprintf("Subscribe to server-side channel %s: (resubscribe: %t, recovered: %t)", e.Channel, e.Resubscribed, e.Recovered))
}

// OnServerUnsubscribe ...
func (h *EventHandler) OnServerUnsubscribe(_ *centrifuge.Client, e centrifuge.ServerUnsubscribeEvent) {
	h.sendMsg2Screen(EventServerUnsubscribe, fmt.Sprintf("Unsubscribe from server-side channel %s", e.Channel))
}

// OnServerJoin ...
func (h *EventHandler) OnServerJoin(_ *centrifuge.Client, e centrifuge.ServerJoinEvent) {
	h.sendMsg2Screen(EventServerJoin, fmt.Sprintf("Server-side join to channel %s: %s (%s)", e.Channel, e.User, e.Client))
}

// OnServerLeave ...
func (h *EventHandler) OnServerLeave(_ *centrifuge.Client, e centrifuge.ServerLeaveEvent) {
	h.sendMsg2Screen(EventServerLeave, fmt.Sprintf("Server-side leave from channel %s: %s (%s)", e.Channel, e.User, e.Client))
}

// OnServerPublish ...
func (h *EventHandler) OnServerPublish(_ *centrifuge.Client, e centrifuge.ServerPublishEvent) {
	h.sendMsg2Screen(EventServerPublish, fmt.Sprintf("Publication from server-side channel %s: %s", e.Channel, e.Data))
}

// OnPublish ...
func (h *EventHandler) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	var chatMessage *ChatMessage
	err := json.Unmarshal(e.Data, &chatMessage)
	if err != nil {
		return
	}
	h.sendMsg2Screen(EventPublish, fmt.Sprintf("Someone says via channel %s: %s", sub.Channel(), chatMessage.Input))
}

// OnJoin ...
func (h *EventHandler) OnJoin(sub *centrifuge.Subscription, e centrifuge.JoinEvent) {
	h.sendMsg2Screen(EventJoin, fmt.Sprintf("Someone joined %s: user id %s, client id %s", sub.Channel(), e.User, e.Client))
}

// OnLeave ...
func (h *EventHandler) OnLeave(sub *centrifuge.Subscription, e centrifuge.LeaveEvent) {
	h.sendMsg2Screen(EventJoin, fmt.Sprintf("Someone left %s: user id %s, client id %s", sub.Channel(), e.User, e.Client))
}

// OnSubscribeSuccess ...
func (h *EventHandler) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	h.sendMsg2Screen(EventSubscribeSuccess, fmt.Sprintf("Subscribed on channel %s, resubscribed: %v, recovered: %v", sub.Channel(), e.Resubscribed, e.Recovered))
}

// OnSubscribeError ...
func (h *EventHandler) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	h.sendMsg2Screen(EventSubscribeError, fmt.Sprintf("Subscribed on channel %s failed, error: %s", sub.Channel(), e.Error))
}

// OnUnsubscribe ...
func (h *EventHandler) OnUnsubscribe(sub *centrifuge.Subscription, _ centrifuge.UnsubscribeEvent) {
	h.sendMsg2Screen(EventUnsubscribe, fmt.Sprintf("Unsubscribed from channel %s", sub.Channel()))
}

func (h *EventHandler) sendMsg2Screen(eventType EventType, msg string) {
	h.gameEventChan <- &Event{
		Type: eventType,
		Msg:  msg,
	}
	log.Print(msg)
}

func (h *EventHandler) newClient() *centrifuge.Client {
	wsURL := "ws://localhost:8000/connection/websocket"
	c := centrifuge.New(wsURL, centrifuge.DefaultConfig())

	// Uncomment to make it work with Centrifugo and its JWT auth.
	c.SetToken(connToken("49", 0))

	c.OnConnect(h)
	c.OnDisconnect(h)
	c.OnMessage(h)
	c.OnError(h)

	c.OnServerPublish(h)
	c.OnServerSubscribe(h)
	c.OnServerUnsubscribe(h)
	c.OnServerJoin(h)
	c.OnServerLeave(h)

	return c
}

// // Init ...
// func (h *EventHandler) Init(eventChan chan *GameEvent) {
// 	h.gameEventChan = eventChan
// }

// Run ...
func (h *EventHandler) Run() {
	defer func() {
		if e := recover(); e != nil {
			h.sendMsg2Screen(EventError, fmt.Sprintf("%v", e))
		}
	}()

	c := h.newClient()
	defer func() { _ = c.Close() }()

	sub, err := c.NewSubscription("news")
	if err != nil {
		log.Panic(err)
	}

	sub.OnPublish(h)
	sub.OnJoin(h)
	sub.OnLeave(h)
	sub.OnSubscribeSuccess(h)
	sub.OnSubscribeError(h)
	sub.OnUnsubscribe(h)

	err = sub.Subscribe()
	if err != nil {
		log.Panic(err)
	}

	pubText := func(text string) error {
		msg := &ChatMessage{
			Input: text,
		}
		data, _ := json.Marshal(msg)
		_, err := sub.Publish(data)
		return err
	}

	err = c.Connect()
	if err != nil {
		log.Panic(err)
	}

	err = pubText("hello")
	if err != nil {
		log.Printf("Error publish: %s", err)
	}

	log.Printf("Print something and press ENTER to send\n")

	// Read input from stdin.
	go func(sub *centrifuge.Subscription) {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			err = pubText(text)
			if err != nil {
				log.Printf("Error publish: %s", err)
			}
		}
	}(sub)

	// Run until CTRL+C.
	select {}
}
