package progimg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"image"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

const defaultPath = "./images"

func init() {
	os.MkdirAll(defaultPath, 0766)
}

func newID() uint64 {
	key := fmt.Sprintf("prog-%d-%v", time.Now().Unix(), rand.Uint64())
	h := fnv.New64()
	h.Write([]byte(key))
	return h.Sum64()
}

func contentTypeOK(ct string) bool {
	for _, t := range []string{"png", "jpeg", "image/png", "image/jpeg"} {
		if ct == t {
			return true
		}
	}

	return false
}

func getPath(id string) string {
	return fmt.Sprintf("%s/%s", defaultPath, id)
}

func saveImage(path string, img *Image) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", path, err)
	}

	defer f.Sync()
	defer f.Close()
	enc := gob.NewEncoder(f)
	return enc.Encode(img)
}

func getImage(path string) (*Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", path, err)
	}

	defer f.Close()
	dec := gob.NewDecoder(f)
	var img Image
	err = dec.Decode(&img)
	return &img, err
}

func getGoImage(img *Image) (image.Image, error) {
	buf := bytes.NewReader(img.Data)
	switch img.Type {
	case "png":
		return png.Decode(buf)
	case "jpeg":
		return jpeg.Decode(buf)
	}

	return nil, fmt.Errorf("unknown image format: %s", img.Type)
}

func transformImage(rct string, img *Image) error {
	if rct == img.Type {
		return nil
	}

	gimg, err := getGoImage(img)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	var buf bytes.Buffer
	switch rct {
	case "png":
		err = png.Encode(&buf, gimg)
	case "jpeg":
		err = jpeg.Encode(&buf, gimg, nil)
	default:
		err = fmt.Errorf("unknown conversion format: %s", rct)
	}

	if err != nil {
		return fmt.Errorf("failed to convert image: %v", err)
	}

	img.Type = rct
	img.Data = buf.Bytes()
	return nil
}
