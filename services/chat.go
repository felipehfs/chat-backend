package services

import (
	"fmt"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
)

// Chat represents
type Chat struct {
	Server *socketio.Server
}

// NewChat .
func NewChat() (*Chat, error) {
	server, err := socketio.NewServer(nil)

	if err != nil {
		return nil, err
	}

	server.OnConnect("connection", func(sock socketio.Conn) error {
		fmt.Println(sock.RemoteHeader())

		return nil
	})

	return &Chat{
		Server: server,
	}, nil
}

func (c Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Server.ServeHTTP(w, r)
}
