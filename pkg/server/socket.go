package server

import (
	"fmt"
	"io"
	"net/http"

	"github.com/cvancleave/go-socket-server/pkg/utils"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
)

func (s *server) handleSocket(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	utils.SetCorsHeaders(w, r)
	websocket.Handler(s.readLoop).ServeHTTP(w, r)
}

func (s *server) readLoop(ws *websocket.Conn) {

	addr := ws.RemoteAddr()
	fmt.Printf("socket - opened: %s\n", addr.String())
	s.connections[ws] = true

	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("error with socket read: %s - %s\n", addr.String(), err.Error())
			continue
		}
		s.broadcast(buf[:n])
	}

	// delete from map
	delete(s.connections, ws)
	fmt.Println("socket - closed:", addr.String())
}

func (s *server) broadcast(bytes []byte) {
	for ws := range s.connections {
		go func(w *websocket.Conn) {
			if _, err := w.Write(bytes); err != nil {
				fmt.Println("socket broadcast error:", err.Error())
			}
		}(ws)
	}
}
