package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/gorilla/websocket"
	"github.com/oitel/tubelas/hub"
)

const (
	readBufferSize  = 1024
	writeBufferSize = 1024
	readLimit       = 512
	pingPeriod      = 50 * time.Second
	pongTimeout     = 60 * time.Second
	writeTimeout    = 10 * time.Second
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrader.Upgrade: ", err)
		return
	}
	defer conn.Close()

	conn.SetReadLimit(512)

	conn.SetReadDeadline(time.Now().Add(pongTimeout))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongTimeout))
		return nil
	})

	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	h := hub.GlobalInstance()
	cl := h.Register()
	defer func() {
		h.Unregister(cl)
		for range cl.Incoming() {
			//
		}
	}()

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("conn.ReadMessage: ", err)
				}
				h.Unregister(cl)
				return
			}

			text := strings.ToValidUTF8(string(msg), string(unicode.ReplacementChar))
			cl.Publish(text)
		}
	}()

	for {
		select {
		case msg, ok := <-cl.Incoming():
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))

			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("json.Marshal: ", err)
				continue
			}
			if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
				log.Println("conn.WriteMessage: ", err)
				return
			}
		case <-ticker.C:
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println("conn.WriteMessage: ", err)
				return
			}
		}
	}
}
