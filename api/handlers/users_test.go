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

func TestGetUser(t *testing.T) {
	req := httptest.NewRequest("GET", "/users/6160fe9757a258c6bdc94056", nil)
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandleUserGet(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	expectedBody := `{"id":"6160fe9757a258c6bdc94056","name":"Souris Ash","email":"sasa@lele.com"}`

	if string(body) != expectedBody {
		t.Errorf("Unexpected body returned. Expected %s and got %s", expectedBody, string(body))
	}
}

func TestCreateUser(t *testing.T) {
	jsonStr := []byte(`{"name":"User N","email":"sasa@lelen.com","password":"thisapass"}`)

	req := httptest.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	senv, client := openDBCon()
	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	senv.HandleUserCreate(w, req)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Handler returned wrong status code: expected %v but received %v.", http.StatusOK, resp.StatusCode)
	}

	if !strings.Contains(string(body), "id") {
		t.Errorf("Body does not contain the ID of inserted object. Body received: %s", string(body))
	}
}
