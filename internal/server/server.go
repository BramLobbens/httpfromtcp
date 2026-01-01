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

func server(address string, options ...func(*Server)) (*Server, error) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, err
	}
	server := &Server{
		state:    atomic.Bool{},
		listener: listener,
	}

	for _, option := range options {
		option(server)
	}

	return server, nil
}

func Serve(port int) (*Server, error) {
	address := net.JoinHostPort("", strconv.Itoa(port))
	server, err := server(address)
	if err != nil {
		return nil, err
	}
	server.state.Store(true)
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
	defer conn.Close()

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/plain\r\n" +
		"Content-Length: 13\r\n" +
		"\r\n" +
		"Hello World!\n"

	if _, err := conn.Write([]byte(response)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}
