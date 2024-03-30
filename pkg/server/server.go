package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cvancleave/go-socket-server/pkg/utils"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/websocket"
)

type server struct {
	port        int
	connections map[*websocket.Conn]bool
}

func Start(port int) {

	// set up server
	s := &server{
		port:        port,
		connections: map[*websocket.Conn]bool{},
	}

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", s.port),
		Handler:     s.routes(),
		IdleTimeout: time.Minute,
	}

	fmt.Println("starting go-socket-server on port:", s.port)

	// serve
	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *server) routes() *httprouter.Router {

	router := httprouter.New()

	router.GET("/socket", s.handleSocket)

	// allow cors for non-get methods
	router.GlobalOPTIONS = http.HandlerFunc(options)

	return router
}

func options(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Access-Control-Request-Method") != "" {
		utils.SetCorsHeaders(w, r)
	}
	w.WriteHeader(http.StatusOK)
}
