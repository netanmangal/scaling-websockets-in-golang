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

// create a new server
// start the websocket server at /ws
// create a handler - whenever there is an incoming request for websocket
// 		accept it
// 		start reading from it.
// whenever something is read, write it to all the connections

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleIncomingWSRequest(ws *websocket.Conn) {
	fmt.Println("New incoming connection from client : ", ws.RemoteAddr())

	// maps in golang are not concurrent safe
	// we should use here some form of mutex
	s.conns[ws] = true
	s.readLoopForSocket(ws)
}

func (s *Server) readLoopForSocket(ws *websocket.Conn) {
	buf := make([]byte, 1024)

	for {
		n, err := ws.Read(buf)

		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection terminated by client : ", ws.RemoteAddr())
				s.conns[ws] = false
				break
			}

			fmt.Println("Read error occured for client : ", ws.RemoteAddr(), " Error : ", err)
			continue
		}

		msg := buf[:n]
		messageToSend := fmt.Sprintf("Message from client : %v is %v", ws.RemoteAddr(), string(msg))
		fmt.Println(messageToSend)

		// propagate the message
		for wsConn, isConnected := range s.conns {
			if isConnected {
				wsConn.Write([]byte(messageToSend))
			}
		}
	}
}

func main() {
	server := NewServer()
	fmt.Println("Setting up the server : ", server)

	http.Handle("/ws", websocket.Handler(server.handleIncomingWSRequest))
	fmt.Println("Setting up the Handler")

	fmt.Println("Server is listening on port 3000")
	http.ListenAndServe(":3000", nil)
}
