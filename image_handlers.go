package progimg

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Image represents an image we store on our end
type Image struct {
	ID     string // id: unique ID for image
	Format string // Format: image format
	Data   []byte // Data: image data
}

// newImage returns a new image from given format and image data
func newImage(ct string, data []byte) *Image {
	return &Image{
		ID:     fmt.Sprint(newID()),
		Format: strings.TrimPrefix(ct, "image/"),
		Data:   data,
	}
}

// uploadTypeHandler aliases function that handles image extraction from request
type uploadTypeHandler func(r *http.Request) (*Image, error)

// uploadTypeHandlers acts a mux for different uploadTypeHandler
var uploadTypeHandlers map[string]uploadTypeHandler

// base64Handler extracts the base64 encoded image from the request
func base64Handler() uploadTypeHandler {
	return uploadTypeHandler(func(r *http.Request) (img *Image, err error) {
		eimg := r.PostForm.Get("image")
		dimg, err := base64.StdEncoding.DecodeString(eimg)
		if err != nil {
			return nil, fmt.Errorf("failed to decode base64 image: %v", err)
		}

		ct := http.DetectContentType(dimg)
		if !contentTypeOK(ct) {
			return nil, fmt.Errorf("unknown content type: %s", ct)
		}

		return newImage(ct, dimg), nil
	})
}

// urlImageHandler fetches the url from request, downloads the image and returns the image
func urlImageHandler() uploadTypeHandler {
	return uploadTypeHandler(func(r *http.Request) (img *Image, err error) {
		iu := r.PostForm.Get("image")
		resp, err := http.Get(iu)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch %s: %v", iu, err)
		}

		defer resp.Body.Close()
		d, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to fecth %s: %v", iu, err)
		}

		ct := resp.Header.Get("Content-type")
		if ct == "" {
			ct = http.DetectContentType(d)
		}

		if !contentTypeOK(ct) {
			return nil, fmt.Errorf("unknown content type found %s: fetch %s", ct, iu)
		}

		return newImage(ct, d), nil
	})
}

// multipartImageHandler extracts the multipart image upload from request
func multipartImageHandler() uploadTypeHandler {
	return uploadTypeHandler(func(r *http.Request) (img *Image, err error) {
		i, _, err := r.FormFile("image")
		if err != nil {
			return nil, fmt.Errorf("failed to fetch multipart image: %v", err)
		}

		defer i.Close()

		d, err := ioutil.ReadAll(i)
		if err != nil {
			return nil, fmt.Errorf("failed to read image file: %v", err)
		}

		ct := http.DetectContentType(d)
		if !contentTypeOK(ct) {
			return nil, fmt.Errorf("unknow content type: %s", ct)
		}

		return newImage(ct, d), nil
	})
}

func init() {
	uploadTypeHandlers = make(map[string]uploadTypeHandler)
	uploadTypeHandlers["base64"] = base64Handler()
	uploadTypeHandlers["url"] = urlImageHandler()
	uploadTypeHandlers["file"] = multipartImageHandler()
}
