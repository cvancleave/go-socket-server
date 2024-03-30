# go-socket-server

A basic implementation of a socket server.

### Running locally

Run server.
- `go run cmd/server/main.go`

Use the Chrome console to connect.
- `let socket = new WebSocket("ws://localhost:4002/socket")`
- `socket.onmessage = (event) => { console.log("received from server: ", event.data) }`
- `socket.send("hello server")`