package progimg

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// handleUpload handles the image upload requests
func handleUpload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	r.ParseMultipartForm(32 << 20)
	imgType := r.FormValue("type")
	h, ok := uploadTypeHandlers[imgType]
	if !ok {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": fmt.Sprintf("unknown format: %s", imgType),
		})
		return
	}

	img, err := h(r)
	if err != nil {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	err = saveImage(getPath(img.ID), img)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	writeJSONResponse(w, http.StatusCreated, map[string]string{
		"id": img.ID,
	})
}

func handleDownload(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		writeJSONResponse(w, http.StatusBadRequest, map[string]string{
			"error": "id is required",
		})
		return
	}
	img, err := getImage(getPath(id))
	if err != nil {
		writeJSONResponse(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	r.ParseForm()
	ct := r.Form.Get("type")
	if ct != "" && ct != img.Type {
		err := transformImage(ct, img)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}

	w.Header().Add("Content-type", fmt.Sprintf("image/%s", img.Type))
	w.WriteHeader(http.StatusOK)
	w.Write(img.Data)
}

func handle404(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	writeJSONResponse(w, http.StatusNotFound, map[string]string{
		"error": "404 not found",
	})
}

func writeJSONResponse(w http.ResponseWriter, c int, d interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(c)
	p, err := json.Marshal(d)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	w.Write(p)
}
