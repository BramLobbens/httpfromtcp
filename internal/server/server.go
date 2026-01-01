package server

import (
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	state    atomic.Bool
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	server := &Server{
		state: atomic.Bool{},
	}
	server.state.Store(true)
	address := net.JoinHostPort("", strconv.Itoa(port))
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server.listener = listener
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.state.Store(false)
	return s.listener.Close()
}

func (s *Server) listen() {
	for s.state.Load() {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.state.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"

	conn.Write([]byte(response))
	err := s.Close()
	if err != nil {
		log.Fatal(err)
	}
}
