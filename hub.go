package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

const (
	defaultBuf = 1000
	serverName = "BLC"
)

type Message struct {
	From    string      `json:"f"`
	To      string      `json:"t"`
	Kind    int         `json:"k"`
	Content interface{} `json:"c"`
	//Sdp       string                      `json:"sdp"`
	//Candidate map[interface{}]interface{} `json:"candidate"`
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {

	// Register requests from the connections.
	register chan *connection
	// Unregister requests from connections.
	unregister chan *connection
	// messages channel
	messages chan []byte
	// broadcast messages from the connections.
	broadcast chan []byte

	// Registered connections.
	connections map[*connection]bool
	//  user->conncetion
	idConnMap map[string]*connection
	idMutex   sync.Mutex
}

var h = hub{
	messages:    make(chan []byte, defaultBuf),
	broadcast:   make(chan []byte, defaultBuf),
	register:    make(chan *connection, defaultBuf),
	unregister:  make(chan *connection, defaultBuf),
	connections: make(map[*connection]bool),
	idConnMap:   make(map[string]*connection),
}

func (h *hub) run() {
	for {
		select {
		case c := <-h.register:
			h.connections[c] = true

			list := make([]string, 0)
			for c := range h.connections {
				list = append(list, c.id)
			}
			listmsg := Message{serverName, "", 3, list}
			b1, _ := json.Marshal(listmsg)
			h.broadcast <- b1
		case c := <-h.unregister:
			if _, ok := h.connections[c]; ok {

				if c.id != "" {
					msg := Message{serverName, "", 0, fmt.Sprintf("%s logout", c.id)}
					b, _ := json.Marshal(msg)
					h.broadcast <- b
				}

				delete(h.connections, c)
				delete(h.idConnMap, c.id)
				close(c.send)

			}
		case b := <-h.messages:
			var m Message
			json.Unmarshal(b, &m)
			if m.Kind == 0 {
				h.broadcast <- b
			} else {
				if user, ok := h.idConnMap[m.To]; ok == true {
					user.send <- b
				}
			}
		case b := <-h.broadcast:
			for c := range h.connections {
				select {
				case c.send <- b:
				default:
					close(c.send)
					delete(h.connections, c)
				}
			}
		}
	}
}

func (h *hub) RegisterUser(in []byte, c *connection) bool {
	var m Message
	json.Unmarshal(in, &m)
	id := m.From
	if id == serverName || id == "*" {
		goto FAIL_RET
	}
	c.id = id

	h.idMutex.Lock()
	defer func() { h.idMutex.Unlock() }()
	if _, ok := h.idConnMap[id]; !ok {
		h.idConnMap[id] = c

		msg := Message{serverName, "", 0, fmt.Sprintf("%s login", id)}
		b, _ := json.Marshal(msg)
		h.broadcast <- b

		list := make([]string, 0)
		for c := range h.connections {
			list = append(list, c.id)
		}
		listmsg := Message{serverName, "", 3, list}
		b1, _ := json.Marshal(listmsg)
		h.broadcast <- b1
		fmt.Println("Register ", id)
		return true
	}

FAIL_RET:
	msg := Message{serverName, "", 0, "login failed: user exsit"}
	b, _ := json.Marshal(msg)
	c.send <- b
	return false
}
