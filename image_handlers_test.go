package progimg

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_newImage(t *testing.T) {
	tests := []struct {
		ct  string
		ect string
	}{
		{
			ct:  "image/png",
			ect: "png",
		},

		{
			ct:  "image/jpeg",
			ect: "jpeg",
		},

		{
			ct:  "pdf",
			ect: "pdf",
		},
	}

	for _, c := range tests {
		img := newImage(c.ct, nil)
		if img.Format != c.ect {
			t.Fatalf("expected %s content type but got %s", c.ect, img.Format)
		}
	}
}

func Test_base64Handler(t *testing.T) {
	tests := []struct {
		image string
		ct    string
		err   string
	}{
		{
			image: getTestBase64("./testdata/testimg.png"),
			ct:    "png",
		},

		{
			image: getTestBase64("./testdata/testimg.jpeg"),
			ct:    "jpeg",
		},

		{
			image: getTestBase64("./testdata/testpdf.pdf"),
			err:   "unknown content type: application/pdf",
		},

		{
			image: "somerandombase64==",
			err:   "failed to decode base64 image",
		},
	}

	b64Handler := base64Handler()
	for _, c := range tests {
		f := url.Values{}
		f.Add("image", c.image)
		r := httptest.NewRequest("POST", "/images", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ParseMultipartForm(32 << 20)
		img, err := b64Handler(r)
		if err != nil {
			if c.err != "" && strings.Contains(err.Error(), c.err) {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if img.Format != c.ct {
			t.Fatalf("expected %s type but got %s", c.ct, img.Format)
		}
	}
}

func Test_urlHandler(t *testing.T) {
	tests := []struct {
		url string
		ct  string
		err string
	}{
		{
			url: "invalid",
			err: "failed to fetch invalid",
		},

		{
			url: "http://che.org.il/wp-content/uploads/2016/12/pdf-sample.pdf",
			err: "unknown content type found application/pdf",
		},

		{
			url: "https://i.vimeocdn.com/portrait/58832_300x300",
			ct:  "jpeg",
		},
	}

	urlHandler := urlImageHandler()
	for _, c := range tests {
		f := url.Values{}
		f.Add("image", c.url)
		r := httptest.NewRequest("POST", "/images", strings.NewReader(f.Encode()))
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.ParseMultipartForm(32 << 20)
		img, err := urlHandler(r)
		if err != nil {
			if c.err != "" && strings.Contains(err.Error(), c.err) {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if img.Format != c.ct {
			t.Fatalf("expected %s type but got %s", c.ct, img.Format)
		}
	}
}

func Test_fileImageUpload(t *testing.T) {
	tests := []struct {
		path string
		ct   string
		err  string
	}{
		{
			path: "./testdata/testpdf.pdf",
			err:  "unknow content type: application/pdf",
		},

		{
			path: "./testdata/testimg.png",
			ct:   "png",
		},
	}

	fileHandler := multipartImageHandler()
	for _, c := range tests {
		file, err := os.Open(c.path)
		if err != nil {
			t.Fatalf("unexpected error: file open: %v", err)
		}

		var body bytes.Buffer
		writer := multipart.NewWriter(&body)
		part, err := writer.CreateFormFile("image", filepath.Base(c.path))
		if err != nil {
			t.Fatalf("unexpected error: multipart create: %v", err)
		}
		_, err = io.Copy(part, file)
		writer.WriteField("type", "file")
		err = writer.Close()
		if err != nil {
			t.Fatalf("unexpected error: multipart writer close: %v", err)
		}

		req := httptest.NewRequest("POST", "/images", &body)
		if err != nil {
			t.Fatalf("unexpected error: post fail: %v", err)
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		file.Close()

		img, err := fileHandler(req)
		if err != nil {
			if c.err != "" && strings.Contains(err.Error(), c.err) {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if img.Format != c.ct {
			t.Fatalf("expected %s type but got %s", c.ct, img.Format)
		}
	}
}
