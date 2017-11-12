package progimg

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleUpload handles the image upload requests
func handleUpload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseForm()
	imgType := r.PostFormValue("type")
	h, ok := uploadTypeHandlers[imgType]
	if !ok {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("unknown format: %s", imgType),
		})
		return
	}

	img, err := h(r)
	if err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = saveImage(getPath(img.ID), img)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	writeResponse(w, http.StatusCreated, map[string]string{
		"id": img.ID,
	})
}

func handle404(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	writeResponse(w, http.StatusNotFound, map[string]string{
		"error": "404 not found",
	})
}

func writeResponse(w http.ResponseWriter, c int, d interface{}) {
	w.WriteHeader(c)
	p, err := json.Marshal(d)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(p)
}
