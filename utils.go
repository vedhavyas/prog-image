package progimg

import (
	"encoding/gob"
	"fmt"
	"hash/fnv"
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
