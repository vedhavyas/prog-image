package progimg

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

// responseWriter holds original response writer and other meta required
type responseWriter struct {
	rw     http.ResponseWriter
	status int
	size   int
}

// Header returns the original response writer's Header
func (w *responseWriter) Header() http.Header {
	return w.rw.Header()
}

// Write write the data to original response writer
func (w *responseWriter) Write(data []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	size, err := w.rw.Write(data)
	w.size += size
	return size, err
}

// WriteHeader writes header to original response writer's WriteHeader
func (w *responseWriter) WriteHeader(header int) {
	w.rw.WriteHeader(header)
	w.status = header
}

//logHandler wraps the handler with logger
func logHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &responseWriter{w, 0, 0}
		handler.ServeHTTP(writer, r)
		end := time.Now()
		latency := end.Sub(start)
		log.Printf("%s [%v] \"%s %s %s\" %d %d \"%s\" %v\n",
			r.RemoteAddr, end.Format(time.RFC1123Z),
			r.Method, r.URL.Path, r.Proto,
			writer.status, writer.size, r.Header.Get("User-Agent"), latency)
	})
}

// recoverHandler defer recover
func recoverHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("[recovered]: %v\n", err)
				fmt.Println(string(debug.Stack()))
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
