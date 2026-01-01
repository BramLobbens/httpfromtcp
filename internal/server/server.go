package server

import (
	"log"
	"net"
	"strconv"
)

type Server struct {
	state    bool // 0 - open, 1 - closed
	listener net.Listener
}

func Serve(port int) (*Server, error) {
	server := &Server{
		state: true, // open
	}
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server.listener = listener
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) listen() {
	for s.state { // while open
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal(err)
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
