package services

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. dMust be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

var wg sync.WaitGroup
var mutex = sync.Mutex{}

// WebsocketClient .
type WebsocketClient struct {
	Upgrader  websocket.Upgrader
	Broadcast chan map[string]interface{}
	Clients   map[*websocket.Conn]bool
	Quit      chan struct{}
}

// NewWebsocketClient .
func NewWebsocketClient() *WebsocketClient {
	broadcast := make(chan map[string]interface{})
	quit := make(chan struct{})
	clients := make(map[*websocket.Conn]bool)

	return &WebsocketClient{
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		Broadcast: broadcast,
		Quit:      quit,
		Clients:   clients,
	}
}

func (wc WebsocketClient) Reader(conn *websocket.Conn) {
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		message := make(map[string]interface{})
		err := conn.ReadJSON(&message)
		if err != nil {
			log.Println("read:", err)

			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					log.Printf("Web socket closed by client: %s", err)
					return
				}
			}
		} else {
			wc.Broadcast <- message
			log.Printf("Sucessfully read %v", message)
		}
	}
	wg.Done()
}

func (wc WebsocketClient) Writer(conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	for {
		select {
		case message := <-wc.Broadcast:
			log.Printf("sending %s", message)

			for client := range wc.Clients {
				mutex.Lock()
				err := client.WriteJSON(message)
				if err != nil {
					log.Println("write:", err)
					delete(wc.Clients, client)
					client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				}
				mutex.Unlock()
			}
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				break
			}
		}
	}

	wg.Done()
	ticker.Stop()
}

func (wc WebsocketClient) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := wc.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	mutex.Lock()
	if _, ok := wc.Clients[c]; !ok {
		wc.Clients[c] = true
	}
	mutex.Unlock()

	defer c.Close()

	wg.Add(2)
	go wc.Reader(c)
	go wc.Writer(c)
	wg.Wait()
}
