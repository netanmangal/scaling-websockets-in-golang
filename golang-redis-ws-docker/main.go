package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/websocket"
)

type Server struct {
	conns       map[*websocket.Conn]bool
	redisClient *redis.Client
}

func NewServer() (s *Server) {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
		redisClient: redis.NewClient(&redis.Options{
			Addr:     os.Getenv("REDISURI"),
			DB:       0,
			Password: "",
		}),
	}
}

func (s *Server) handleIncomingWSRequest(ctx context.Context, ws *websocket.Conn) {
	fmt.Println("New incoming ws connection request from client : ", ws.RemoteAddr())
	s.conns[ws] = true
	s.readLoopForSocket(ctx, ws)
}

func (s *Server) readLoopForSocket(ctx context.Context, ws *websocket.Conn) {
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
		s.redisClient.Publish(ctx, "message-chat", msg)
	}
}

func (s *Server) broadcast(msg []byte) {
	for k, v := range s.conns {
		if v {
			k.Write(msg)
		}
	}
}

func (s *Server) SubscribeRedis(ctx context.Context) {
	sub := s.redisClient.Subscribe(ctx, "message-chat")
	defer sub.Close()

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			fmt.Println("Error receiving messages from Redis")
		}

		if msg.Payload != "" {
			s.broadcast([]byte(msg.Payload))
			fmt.Println("Broadcasting complete for message : ", msg.Payload)
		}
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading envs : ", err)
	}

	server := NewServer()
	ctx := context.Background()
	go server.SubscribeRedis(ctx)

	http.Handle("/ws", websocket.Handler(func(conn *websocket.Conn) {
		server.handleIncomingWSRequest(ctx, conn)
	}))
	fmt.Println("Setting up the handler")

	fmt.Println("Server is listening on port : ", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		fmt.Println("Error starting the server : ", err)
	}
}
