package progimg

import (
	"net/http/httptest"
	"net/url"
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
		if img.Type != c.ect {
			t.Fatalf("expected %s content type but got %s", c.ect, img.Type)
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

		if img.Type != c.ct {
			t.Fatalf("expected %s type but got %s", c.ct, img.Type)
		}
	}
}
