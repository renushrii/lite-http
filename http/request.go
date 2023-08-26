package http

import "io"

type Request struct {
	Method string
	Path   string

	version string

	Headers map[string]string
	Cookies map[string]string

	body io.ReadCloser
}

func (req *Request) Body() io.ReadCloser {
	return req.body
}
