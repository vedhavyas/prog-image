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

func getTestBase64(path string) string {
	f, err := os.Open(path)
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
	form.Add("image", getTestBase64("./testdata/testimg.png"))
	req, _ := http.NewRequest("POST", s.URL+"/images", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error: response body: %v", err)
	}

	var res struct {
		ID string
	}

	err = json.Unmarshal(d, &res)
	if err != nil {
		t.Fatalf("unexpected error: json marshalling : %v", err)
	}

	return res.ID
}

func Test_uploadImage_base64(t *testing.T) {
	s := setup()
	postTestImage(t, s)
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
		t.Fatalf("unexpected error: post response: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	cleanup(s)
}

func multipartTestRequest(t *testing.T, s *httptest.Server, path string) (req *http.Request) {
	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("unexpected error: file open: %v", err)
	}
	defer file.Close()

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("image", filepath.Base(path))
	if err != nil {
		t.Fatalf("unexpected error: multipart create: %v", err)
	}
	_, err = io.Copy(part, file)
	writer.WriteField("type", "file")
	err = writer.Close()
	if err != nil {
		t.Fatalf("unexpected error: multipart writer close: %v", err)
	}

	req, err = http.NewRequest("POST", s.URL+"/images", &body)
	if err != nil {
		t.Fatalf("unexpected error: post fail: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req
}

func Test_uploadImageFile(t *testing.T) {
	s := setup()
	path := "./testdata/testimg.png"
	req := multipartTestRequest(t, s, path)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: get fail: %v", err)
	}

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	cleanup(s)
}

func Test_unknownType(t *testing.T) {
	s := setup()
	form := url.Values{}
	form.Add("type", "random")
	form.Add("image", getTestBase64("./testdata/testimg.png"))
	req, _ := http.NewRequest("POST", s.URL+"/images", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	cleanup(s)
}

func Test_uploadBase64_error(t *testing.T) {
	s := setup()
	form := url.Values{}
	form.Add("type", "base64")
	form.Add("image", getTestBase64("./testdata/testpdf.pdf"))
	req, _ := http.NewRequest("POST", s.URL+"/images", strings.NewReader(form.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	cleanup(s)
}

func Test_uploadURL_error(t *testing.T) {
	s := setup()
	form := url.Values{}
	form.Add("type", "url")
	form.Add("image", "http://che.org.il/wp-content/uploads/2016/12/pdf-sample.pdf")
	resp, err := http.PostForm(s.URL+"/images/", form)
	if err != nil {
		t.Fatalf("unexpected error: post fail: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	cleanup(s)
}

func Test_uploadImageFile_error(t *testing.T) {
	s := setup()
	path := "./testdata/testpdf.pdf"
	req := multipartTestRequest(t, s, path)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("unexpected error: post fail: %v", err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
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
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "image/png" {
		t.Fatalf("unexpected error: content-type: %s", ct)
	}

	rd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error: response body: %v", err)
	}

	if getTestBase64("./testdata/testimg.png") != base64.StdEncoding.EncodeToString(rd) {
		t.Fatalf("unexpected error: wrong image: %v", rd)
	}

	cleanup(s)
}

func Test_downloadImage_convert_PNG_JPEG(t *testing.T) {
	s := setup()
	id := postTestImage(t, s)
	resp, err := http.Get(s.URL + "/images/" + id + "?type=jpeg")
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if ct != "image/jpeg" {
		t.Fatalf("unexpected error: content-type: %s", ct)
	}

	rd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("unexpected error: response body: %v", err)
	}

	if getTestBase64("./testdata/testimg.jpeg") != base64.StdEncoding.EncodeToString(rd) {
		t.Fatalf("unexpected error: wrong image: %v", rd)
	}

	cleanup(s)
}

func Test_downloadImage_covert_unknown(t *testing.T) {
	s := setup()
	id := postTestImage(t, s)
	resp, err := http.Get(s.URL + "/images/" + id + "?type=pdf")
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected error: status code: %d", resp.StatusCode)
	}
}
