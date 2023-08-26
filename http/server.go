package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type HandlerFunc func(Request, *Response) error

type Server struct {
	addr     string
	handlers map[string]HandlerFunc
}

func NewServer(addr string) *Server {
	return &Server{
		addr:     addr,
		handlers: make(map[string]HandlerFunc),
	}
}

func (s *Server) Get(path string, handler HandlerFunc) {
	key := fmt.Sprintf("GET:%s", path)
	s.handlers[key] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	fmt.Printf("Server is listening on %s\n", s.addr)
	fmt.Println("Registered paths:", s.handlers)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	buf := bufio.NewReader(conn)
	// fmt.Println(buf)
	line, _, err := buf.ReadLine()

	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	lineParts := strings.Split(string(line), " ")
	if len(lineParts) != 3 {
		fmt.Println("Malformed request line")
		return
	}

	method := lineParts[0]
	path := lineParts[1]

	headers := make(map[string]string)
	cookies := make(map[string]string)

	for {
		headerLine, _, err := buf.ReadLine()
		if err != nil || len(headerLine) == 0 {
			break
		}
		parts := strings.SplitN(string(headerLine), ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "Cookie" {
			cookieParts := strings.Split(value, ";")
			for _, cookiePart := range cookieParts {
				cookie := strings.SplitN(strings.TrimSpace(cookiePart), "=", 2)
				if len(cookie) == 2 {
					cookies[cookie[0]] = cookie[1]
				}
			}
		} else {
			headers[key] = value
		}
	}

	contentLength, _ := strconv.Atoi(headers["Content-Length"])
	// fmt.Println(contentLength)

	req := Request{
		Method:  method,
		Path:    path,
		Headers: headers,
		Cookies: cookies,
		body:    io.NopCloser(io.LimitReader(buf, int64(contentLength))),
	}

	resp := NewResponse(conn)
	defer resp.flush()

	hkey := fmt.Sprintf("%s:%s", strings.ToUpper(method), path)
	fmt.Println(hkey)
	handler, ok := s.handlers[hkey]
	if !ok {
		resp.StatusCode(http.StatusNotFound)
		resp.flush()
		return
	}

	// fmt.Println("Registered paths:", s.handlers)

	err = handler(req, resp)
	if err != nil {
		fmt.Println("Error handling request:", err)
	}
}
