package kaca

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
	SUB_PREFIX     = "__sub:"
	PUB_PREFIX     = "__pub:"
	maxTopics      = 100
	SPLIT_LINE     = "_:_"
)

var disp = NewDispatcher()

type connection struct {
	// websocket connection.
	ws     *websocket.Conn
	send   chan []byte
	topics []string
	cid    uint64
}

func (c *connection) deliver() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		if err := c.ws.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if err := c.sendMsg(websocket.CloseMessage, []byte{}); err != nil {
					fmt.Println(err)
				}
				return
			}
			if err := c.sendMsg(websocket.TextMessage, message); err != nil {
				fmt.Println(err)
				return
			}
		case <-ticker.C:
			if err := c.sendMsg(websocket.PingMessage, []byte{}); err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func (c *connection) dispatch() {
	defer func() {
		disp.unregister <- c
		if err := c.ws.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		fmt.Println(err)
	}
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetPongHandler(func(string) error {
		if err := c.ws.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	})
	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		msg := string(message)
		if strings.Contains(msg, SUB_PREFIX) {
			topic := strings.Split(msg, SUB_PREFIX)[1]
			disp.sub <- strconv.Itoa(int(c.cid)) + SPLIT_LINE + topic
		} else if strings.Contains(msg, PUB_PREFIX) {
			topic_msg := strings.Split(msg, PUB_PREFIX)[1]
			disp.pub <- topic_msg
		} else {
			disp.broadcast <- message
		}
	}
}

func (c *connection) sendMsg(mt int, payload []byte) error {
	if err := c.ws.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		fmt.Println(err)
		return err
	}
	return c.ws.WriteMessage(mt, payload)
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{cid: uint64(rand.Int63()), send: make(chan []byte, 256), ws: ws, topics: make([]string, maxTopics)}
	disp.register <- c
	go c.dispatch()
	c.deliver()
}

func serveWsCheckOrigin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	c := &connection{cid: uint64(rand.Int63()), send: make(chan []byte, 256), ws: ws, topics: make([]string, maxTopics)}
	disp.register <- c
	go c.dispatch()
	c.deliver()
}
