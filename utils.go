package progimg

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"os"
	"time"
)

// defaultPath to store the images
const defaultPath = "./images"

// whiteBackground while converting from png to jpeg
var whiteBackground = color.RGBA{0xff, 0xff, 0xff, 0xff}

// supportedCTs holds all the supported content types
var supportedCTs = []string{"png", "jpeg", "image/png", "image/jpeg"}

func init() {
	os.MkdirAll(defaultPath, 0766)
}

// newID returns a new unique id
func newID() uint64 {
	key := fmt.Sprintf("prog-%d-%v", time.Now().Unix(), rand.Uint64())
	h := fnv.New64()
	h.Write([]byte(key))
	return h.Sum64()
}

// contentTypeOK checks if the given content-type present in the supported list
func contentTypeOK(ct string) bool {
	for _, t := range supportedCTs {
		if ct == t {
			return true
		}
	}

	return false
}

// getPath constructs the image path
func getPath(id string) string {
	return fmt.Sprintf("%s/%s", defaultPath, id)
}

// saveImage will save the image at given path using gob encoding
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

// getImage will extract the image from the file using gob decoder
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

// getGoImage returns image.Image from our Image
func getGoImage(img *Image) (image.Image, error) {
	buf := bytes.NewReader(img.Data)
	switch img.Format {
	case "png":
		return png.Decode(buf)
	case "jpeg":
		return jpeg.Decode(buf)
	}

	return nil, fmt.Errorf("unknown image format: %s", img.Format)
}

// transformImage will transform image to rct format
func transformImage(rct string, img *Image) error {
	if rct == img.Format {
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
		dst := image.NewRGBA(gimg.Bounds())
		draw.Draw(dst, dst.Bounds(), image.NewUniform(whiteBackground),
			image.Point{}, draw.Src)
		draw.Draw(dst, dst.Bounds(), gimg, gimg.Bounds().Min, draw.Over)
		err = jpeg.Encode(&buf, dst, nil)
	default:
		err = fmt.Errorf("unknown conversion format: %s", rct)
	}

	if err != nil {
		return fmt.Errorf("failed to convert image: %v", err)
	}

	img.Format = rct
	img.Data = buf.Bytes()
	return nil
}
