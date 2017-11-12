package progimg

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

type Image struct {
	ID   string
	Type string
	Data []byte
}

func newImage(ct string, data []byte) *Image {
	return &Image{
		ID:   fmt.Sprint(newID()),
		Type: ct,
		Data: data,
	}
}

type uploadTypeHandler func(r *http.Request) (*Image, error)

var uploadTypeHandlers map[string]uploadTypeHandler

func base64Handler() uploadTypeHandler {
	return uploadTypeHandler(func(r *http.Request) (img *Image, err error) {
		ct := r.PostForm.Get("content-type")
		if !contentTypeOK(ct) {
			return nil, fmt.Errorf("unknown content type: %s", ct)
		}

		eimg := r.PostForm.Get("image")
		dimg, err := base64.StdEncoding.DecodeString(eimg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image: %v", err)
		}

		return newImage(ct, dimg), nil
	})
}

func init() {
	uploadTypeHandlers = make(map[string]uploadTypeHandler)
	uploadTypeHandlers["base64"] = base64Handler()
}
