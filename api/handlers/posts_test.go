package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetPost(t *testing.T) {
	req := httptest.NewRequest("GET", "/posts/6161578d7ca34c010e0f21d8", nil)
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandlePostGet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	expectedBody := `{"id":"6161578d7ca34c010e0f21d8","posted_by":"616156d49ab2934adcee255e","caption":"Another caption","img_url":"some.url.here","posted_on":"2021-10-09T08:49:17.482Z"}`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}
}

func TestGetPostBadID(t *testing.T) {
	req := httptest.NewRequest("GET", "/posts/616157e0f21d8", nil)
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandlePostGet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	expectedBody := "Bad postID"

	if strings.TrimRight(string(body), "\n") != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, strings.TrimRight(string(body), "\n"))
	}
}

func TestGetNonExistentPostID(t *testing.T) {
	req := httptest.NewRequest("GET", "/posts/6160578d7ca34c010e0f21d8", nil)
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandlePostGet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	expectedBody := `{}`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}
}

func TestCreatePost(t *testing.T) {
	jsonStr := []byte(`{"posted_by":"616156d49ab2934adcee255e","caption":"Caption 14","img_url":"sample.url.here"}`)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandlePostCreate(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	if !strings.Contains(string(body), "id") {
		t.Errorf("Body does not contain the ID of inserted object. Body received: %s", string(body))
	}
}

func TestCreatePostEmptyFields(t *testing.T) {
	jsonStr := []byte(`{"posted_by":"","caption":"","img_url":"sample.url.here"}`)

	req := httptest.NewRequest("POST", "/posts", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandlePostCreate(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if strings.TrimRight(string(body), "\n") != "Bad Request" {
		t.Errorf("Expected Bad Request in body. Body received: %s", string(body))
	}
}
