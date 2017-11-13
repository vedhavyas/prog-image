package progimg

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware_Recover(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("panicking...")
	})

	r, _ := http.NewRequest("GET", "http://localhost:4040", nil)
	w := httptest.NewRecorder()
	func() {
		defer func() {
			err := recover()
			if err == nil {
				t.Fatalf("Expected to panic but didn't\n")
			}
		}()
		h.ServeHTTP(w, r)
	}()

	rh := recoverHandler(h)
	rh.ServeHTTP(w, r)
}
