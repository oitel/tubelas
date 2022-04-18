package web

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
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

	msgs := make(chan []byte)

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("conn.ReadMessage: ", err)
				}
				close(msgs)
				return
			}
			msgs <- msg
		}
	}()

	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
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
