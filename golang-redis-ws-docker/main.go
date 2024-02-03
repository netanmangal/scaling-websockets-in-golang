package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

func NewServer() (s *Server) {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleIncomingWSRequest(ws *websocket.Conn) {
	fmt.Println("New incoming ws connection request from client : ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoopForSocket(ws)
}

func (s *Server) readLoopForSocket(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)
		if err != nil {
			fmt.Println("Error reading from WS client : ", ws.RemoteAddr(), " - error : ", err)
			break
		} else if err == io.EOF {
			fmt.Println("WS Disconnected : ", ws.RemoteAddr())
			s.conns[ws] = false
			break
		}

		msg := buf[:n]
		fmt.Println("message : ", string(msg), " - client: ", ws.RemoteAddr())
		s.broadcast(msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	for k, v := range s.conns {
		if v {
			k.Write(msg)
		}
	}
}

func main() {
	server := NewServer()

	http.Handle("/ws", websocket.Handler(server.handleIncomingWSRequest))
	fmt.Println("Setting up the handler")

	fmt.Println("Server is listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
