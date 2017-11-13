package progimg

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Image struct {
	ID   string
	Type string
	Data []byte
}

func newImage(ct string, data []byte) *Image {
	return &Image{
		ID:   fmt.Sprint(newID()),
		Type: strings.TrimPrefix(ct, "image/"),
		Data: data,
	}
}

type uploadTypeHandler func(r *http.Request) (*Image, error)

var uploadTypeHandlers map[string]uploadTypeHandler

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

func multipartImageHandler() uploadTypeHandler {
	return uploadTypeHandler(func(r *http.Request) (img *Image, err error) {
		r.ParseMultipartForm(32 << 20)
		i, _, err := r.FormFile("image")
		if err != nil {
			log.Println("read error")
			return nil, fmt.Errorf("failed to fetch multipart image: %v", err)
		}

		defer i.Close()

		d, err := ioutil.ReadAll(i)
		if err != nil {
			log.Println("read error")
			return nil, fmt.Errorf("failed to read image file: %v", err)
		}

		ct := http.DetectContentType(d)
		if !contentTypeOK(ct) {
			log.Println("content type:", ct)
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
