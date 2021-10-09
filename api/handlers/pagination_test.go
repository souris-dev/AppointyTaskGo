package handlers

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserPostsGet(t *testing.T) {
	// pagination: first request

	firstGetRequestBody := []byte(`{"last_id":"6161882093c27946c57c996a","last_posted_on":"2021-10-09T12:16:32.361Z","n_new":3,"first_request":true}`)
	req := httptest.NewRequest("GET", "/posts/users/616156d49ab2934adcee255e", bytes.NewBuffer(firstGetRequestBody))
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandleUserPostsGet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	expectedBody := `[{"id":"6161578d7ca34c010e0f21d8","posted_by":"616156d49ab2934adcee255e","caption":"Another caption","img_url":"some.url.here","posted_on":"2021-10-09T08:49:17.482Z"},{"id":"6161872d93c27946c57c9969","posted_by":"616156d49ab2934adcee255e","caption":"Caption 3","img_url":"some.url.here3","posted_on":"2021-10-09T12:12:29.838Z"},{"id":"6161882093c27946c57c996a","posted_by":"616156d49ab2934adcee255e","caption":"Caption 4","img_url":"some.url.here4","posted_on":"2021-10-09T12:16:32.361Z"}]`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}

	// pagination: next request

	secondGetRequestBody := []byte(`{"last_id":"6161882093c27946c57c996a","last_posted_on":"2021-10-09T12:16:32.361Z","n_new":3,"first_request":false}`)

	req = httptest.NewRequest("GET", "/posts/users/616156d49ab2934adcee255e", bytes.NewBuffer(secondGetRequestBody))
	w = httptest.NewRecorder()

	senv.HandleUserPostsGet(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	expectedBody = `[{"id":"6161883493c27946c57c996b","posted_by":"616156d49ab2934adcee255e","caption":"Caption 5","img_url":"some.url.here5","posted_on":"2021-10-09T12:16:52.558Z"},{"id":"6161884393c27946c57c996c","posted_by":"616156d49ab2934adcee255e","caption":"Caption 6","img_url":"some.url.here6","posted_on":"2021-10-09T12:17:07.665Z"},{"id":"6161884793c27946c57c996d","posted_by":"616156d49ab2934adcee255e","caption":"Caption 7","img_url":"some.url.here6","posted_on":"2021-10-09T12:17:11.478Z"}]`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}

	// pagination: third request

	thirdGetRequestBody := []byte(`{"last_id":"6161884793c27946c57c996d","last_posted_on":"2021-10-09T12:17:11.478Z","n_new":3}`)

	req = httptest.NewRequest("GET", "/posts/users/616156d49ab2934adcee255e", bytes.NewBuffer(thirdGetRequestBody))
	w = httptest.NewRecorder()

	senv.HandleUserPostsGet(w, req)

	resp = w.Result()
	body, _ = ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if err := checkResponseHeaders(resp); err != nil {
		t.Errorf(err.Error())
	}

	expectedBody = `[{"id":"6161884f93c27946c57c996e","posted_by":"616156d49ab2934adcee255e","caption":"Caption 8","img_url":"some.url.here6","posted_on":"2021-10-09T12:17:19.805Z"},{"id":"6161885793c27946c57c996f","posted_by":"616156d49ab2934adcee255e","caption":"Caption 9","img_url":"some.url.here6","posted_on":"2021-10-09T12:17:27.653Z"},{"id":"6161885f93c27946c57c9970","posted_by":"616156d49ab2934adcee255e","caption":"Caption 10","img_url":"some.url.here6","posted_on":"2021-10-09T12:17:35.188Z"}]`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}
}
