package http

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	httpVersion      string
	statusCode       int
	statusCodeString string

	cookies map[string]string
	headers map[string]string
	conn    net.Conn

	bodyBuffer *bytes.Buffer // Buffer to hold the response body
}

func NewResponse(conn net.Conn) *Response {
	return &Response{
		conn:       conn,
		headers:    make(map[string]string),
		cookies:    make(map[string]string),
		bodyBuffer: bytes.NewBuffer(nil),
	}
}

func (resp *Response) Write(data []byte) (int, error) {
	return resp.bodyBuffer.Write(data)
}

func (resp *Response) StatusCode(statusCode int) {
	resp.statusCode = statusCode
	resp.statusCodeString = http.StatusText(statusCode)
}

func (resp *Response) AddHeader(key, value string) {
	resp.headers[key] = value
}

func (resp *Response) SetCookie(key, value string) {
	resp.cookies[key] = value
}

// flush sends the response to the client
func (resp *Response) flush() error {
	respBytes := resp.bodyBuffer.Bytes()
	size := len(respBytes)

	// Write the response line and headers
	resp.conn.Write([]byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", resp.statusCode, resp.statusCodeString)))
	for key, value := range resp.headers {
		resp.conn.Write([]byte(fmt.Sprintf("%s: %s\r\n", key, value)))
	}
	for key, value := range resp.cookies {
		resp.conn.Write([]byte(fmt.Sprintf("Set-Cookie: %s=%s\r\n", key, value)))
	}

	resp.conn.Write([]byte(fmt.Sprintf("Content-Length: %d\r\n", size)))
	resp.conn.Write([]byte("\r\n")) // End of headers

	// Write the response body
	resp.conn.Write(respBytes)

	return nil
}
