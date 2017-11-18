package mulekick

// fake response writer wrapper to tell us when middleware writes
import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"
)

type ResponseWriter struct {
	ResponseWriter  http.ResponseWriter
	responseWritten bool
	statusCode      int
	start           time.Time
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w, false, http.StatusOK, time.Now()}
}

func (wr *ResponseWriter) Header() http.Header {
	return wr.ResponseWriter.Header()
}

func (wr *ResponseWriter) Write(b []byte) (int, error) {
	wr.responseWritten = true
	return wr.ResponseWriter.Write(b)
}

func (wr *ResponseWriter) WriteHeader(header int) {
	wr.responseWritten = true
	wr.statusCode = header
	wr.ResponseWriter.WriteHeader(header)
}

func (wr *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	wr.responseWritten = true

	if hijacker, ok := wr.ResponseWriter.(http.Hijacker); ok {
		return hijacker.Hijack()
	}

	return nil, nil, errors.New("mulekick: response does not implement http.Hijacker")
}
