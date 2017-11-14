package progimg

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_contentTypeOK(t *testing.T) {
	tests := []struct {
		ct string
		r  bool
	}{
		{
			ct: "png",
			r:  true,
		},
		{
			ct: "gif",
			r:  false,
		},
		{
			ct: "pdf",
			r:  false,
		},
		{
			ct: "jpeg",
			r:  true,
		},
	}

	for _, c := range tests {
		r := contentTypeOK(c.ct)
		if r != c.r {
			t.Fatalf("Unexpected error: %s, %v", c.ct, r)
		}
	}
}

func Test_saveImage_getImage(t *testing.T) {
	tests := []*Image{
		{
			ID:     "12345",
			Format: "image/png",
			Data:   []byte{1, 200, 32, 23},
		},

		{
			ID:     "dsbfdbhd",
			Format: "image/jpeg",
			Data:   []byte{1, 23, 12, 123},
		},
	}

	path := "./images"
	os.Mkdir(path, 0766)
	for _, i := range tests {
		p := fmt.Sprintf("%s/%s", path, i.ID)
		err := saveImage(p, i)
		if err != nil {
			t.Fatalf("unexpected error: save image: %v", err)
		}

		img, err := getImage(p)
		if err != nil {
			t.Fatalf("unexpected error: get image: %v", err)
		}

		if !reflect.DeepEqual(i, img) {
			t.Fatal("unexpected error: image mismatch")
		}
	}
}

func Test_transformImage(t *testing.T) {
	tests := []struct {
		ct   string
		data string
		rct  string
		tfp  string
		err  string
	}{
		{
			ct:   "png",
			data: getTestBase64("./testdata/testimg.png"),
			tfp:  "./testdata/testimg.jpeg",
			rct:  "jpeg",
		},

		{
			ct:   "png",
			data: getTestBase64("./testdata/testimg.png"),
			tfp:  "./testdata/testimg.png",
			rct:  "png",
		},

		{
			ct:   "jpeg",
			data: getTestBase64("./testdata/testimg.jpeg"),
			tfp:  "./testdata/testjpegtopng.png",
			rct:  "png",
		},

		{
			ct:   "png",
			data: getTestBase64("./testdata/testimg.png"),
			rct:  "pdf",
			err:  "unknown conversion format: pdf",
		},
	}

	for _, c := range tests {
		data, _ := base64.StdEncoding.DecodeString(c.data)
		img := newImage(c.ct, data)
		err := transformImage(c.rct, img)
		if err != nil {
			if strings.Contains(err.Error(), c.err) {
				continue
			}

			t.Fatalf("unexpected error: %v", err)
		}

		if img.Format != c.rct {
			t.Fatalf("format mismatch: %s != %s", c.rct, img.Format)
		}

		edata, _ := base64.StdEncoding.DecodeString(getTestBase64(c.tfp))
		if !bytes.Equal(edata, img.Data) {
			t.Fatalf("image data mismatch: %v != %v", edata, img.Data)
		}
	}
}
