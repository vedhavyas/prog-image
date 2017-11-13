package progimg

import (
	"fmt"
	"os"
	"reflect"
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
			ID:   "12345",
			Type: "image/png",
			Data: []byte{1, 200, 32, 23},
		},

		{
			ID:   "dsbfdbhd",
			Type: "image/jpeg",
			Data: []byte{1, 23, 12, 123},
		},
	}

	path := "./images-test"
	os.Mkdir(path, 0777)
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
