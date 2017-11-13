package progimg

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func setup() *httptest.Server {
	r := getRouter()
	return httptest.NewServer(r)
}

func cleanup(s *httptest.Server) {
	s.Close()
}

func getTestImgBase64() string {
	f, err := os.Open("./testdata/testimg.png")
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	return base64.StdEncoding.EncodeToString(d)
}

func postTestImage(t *testing.T, s *httptest.Server) (id string) {
	form := url.Values{}
	form.Add("type", "base64")
	form.Add("content-type", "png")
	form.Add("image", getTestImgBase64())
	req, _ := http.NewRequest("POST", s.URL+"/images", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: %v", resp)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var res struct {
		ID string
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}

	return res.ID
}

func Test_uploadImage_base64(t *testing.T) {
	s := setup()
	form := url.Values{}
	form.Add("type", "base64")
	form.Add("image", getTestImgBase64())
	req, _ := http.NewRequest("POST", s.URL+"/images", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: %v", resp)
	}

	cleanup(s)
}

func Test_downloadImage(t *testing.T) {
	s := setup()
	id := postTestImage(t, s)
	resp, err := http.Get(s.URL + "/images/" + id)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected error: %v", resp)
	}

	rd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if getTestImgBase64() != base64.StdEncoding.EncodeToString(rd) {
		t.Fatalf("unexpected error: wrong image: %v", rd)
	}

	cleanup(s)
}

func Test_uploadImageURL(t *testing.T) {
	s := setup()
	id := postTestImage(t, s)
	u := s.URL + "/images/" + id
	form := url.Values{}
	form.Add("type", "url")
	form.Add("image", u)
	resp, err := http.PostForm(s.URL+"/images/", form)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: %v", resp)
	}

	cleanup(s)
}

func Test_uploadImageFile(t *testing.T) {
	s := setup()
	path := "./testdata/testimg.png"
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = io.Copy(part, file)
	writer.WriteField("type", "file")
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req, err := http.NewRequest("POST", s.URL+"/images", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code: %v", resp)
	}

	cleanup(s)
}
