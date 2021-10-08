package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	UserID   uint32
	Name     string
	Email    string
	Password string
}

type Post struct {
	PostID   uint64
	Caption  string
	ImgURL   string
	PostedOn time.Time
}

// utilities

func getHashed256(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", hash)
}

func makeCheckMethodHandler(method string, handlerFn func(writer http.ResponseWriter, req *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, fmt.Sprintf("This endpoint only accepts %s requests!", method), http.StatusBadRequest)
			return
		} else {
			handlerFn(w, req)
		}
	}
}

// handlers

func handleUserCreate(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello from handleUserCreate! Invoked at %s.", req.URL.Path[:])

	var user User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(writer, "Creating user: %v", user)
}

func handleUserGet(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello from handleUserGet as a %s request! Invoked at %s.", req.Method, req.URL.Path[:])

	pass := "just a pass"
	user := User{
		0,
		"Souris Ash",
		"sasa@lele.com",
		getHashed256(pass),
	}

	jsonPost, err := json.Marshal(user)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(writer, string(jsonPost))
}

func handlePostCreate(writer http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(writer, "Hello from handlePostCreate! Invoked at %s.", req.URL.Path[:])

	var post Post

	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(writer, "Creating post: %v", post)
}

func handlePostGet(writer http.ResponseWriter, req *http.Request) {
	post := Post{0, "Hello from post 0!", "this is an url", time.Now()}
	jsonPost, err := json.Marshal(post)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(writer, "Hello from handlePostGet! Invoked at %s.\n", req.URL.Path[:])
	fmt.Fprintf(writer, string(jsonPost))
}

func handleUserPostsGet(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello from handleUserPostsGet! Invoked at %s.", req.URL.Path[:])
}

func main() {
	http.HandleFunc("/users", makeCheckMethodHandler("POST", handleUserCreate))
	http.HandleFunc("/users/", makeCheckMethodHandler("GET", handleUserGet))
	http.HandleFunc("/posts", makeCheckMethodHandler("POST", handlePostCreate))
	http.HandleFunc("/posts/", makeCheckMethodHandler("GET", handlePostGet))
	http.HandleFunc("/posts/users/", makeCheckMethodHandler("GET", handleUserPostsGet))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
