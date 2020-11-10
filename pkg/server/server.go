package server

import (
	"strconv"
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HandlerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HandlerFunc)}
}

func (s *Server) Register(path string, handler HandlerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	log.Println("Server is listening!")
	for {
		conn, err := listener.Accept()
		if err != nil {
			conn.Close()
			continue
		}

		log.Println("All set - you are connected!")

		go s.handle(conn)
	}
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, (4096))

	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
		err = nil
	}
	if err != nil {
		return
	}
	data := buf[:n]
	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)

	if requestLineEnd == -1 {
		log.Print("Error with the request line")
		return
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")

	path := parts[1]

	s.mu.RLock()
	handler, ok := s.handlers[path]
	s.mu.RUnlock()
	if !ok {
		return
	}
	handler(conn)
	return
}

func (s *Server) Requesting(body string) func(conn net.Conn){
	return func(conn net.Conn) {
		_, err := conn.Write([]byte(
			"HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"Content-Type: text/html\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			body,
		))
		if err != nil {
			log.Print(err)
		}
	}
}