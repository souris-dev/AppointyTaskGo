package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ServerEnv struct {
	DB *mongo.Database
}

// the struct tags (other than the first one) below are not really required
// but kept here to ease connection between the frontend and backend
// in case the fields in JSON have a different name

type User struct {
	// JSON unmarshalling should skip this field
	// UserID is an objectID provided by mongodb on insertion
	UserID string `json:"id,omitempty" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email" bson:"email"`

	// The following field initially contains the original password just after
	// JSON unmarshalling of the request body, after which it is hashed and updated for storage.
	// Since we do not have a control over the frontend, we're assuming
	// that we are getting the password as plaintext.
	PwdHash string `json:"password,omitempty" bson:"p_hash"`
}

// This struct defines the form in which the User data is returned in
// the GET endpoint. The password is not returned.

type UserReturn struct {
	UserID string `json:"id" bson:"_id"`
	Name   string `json:"name" bson:"name"`
	Email  string `json:"email" bson:"email"`
}

type Post struct {
	PostID      string    `json:"id,omitempty" bson:"_id"`
	PostedByUID string    `json:"posted_by" bson:"posted_by"`
	Caption     string    `json:"caption" bson:"caption"`
	ImgURL      string    `json:"img_url" bson:"img_url"`
	PostedOn    time.Time `json:"posted_on,omitempty" bson:"posted_on"` // filled at the server
}

// utilities

func getHashed256(pass string) string {
	hash := sha256.Sum256([]byte(pass))
	return fmt.Sprintf("%x", hash)
}

// This function is a wrapper for checking the correct HTTP verb is used
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

// Function to add appropriate headers
func addHeadersUtil(headers map[string]string, w *http.ResponseWriter) {
	for key, val := range headers {
		(*w).Header().Set(key, val)
	}
}

// Function to add common headers
func addCommonHeaders(w *http.ResponseWriter) {
	addHeadersUtil(map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}, w)
}

// handlers

func (senv *ServerEnv) handleUserCreate(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello from handleUserCreate! Invoked at %s.", req.URL.Path[:])

	var user User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	addCommonHeaders(&writer)
	fmt.Fprintf(writer, "Creating user: %v", user)
}

func (senv *ServerEnv) handleUserGet(writer http.ResponseWriter, req *http.Request) {
	urlParts := strings.Split(req.URL.Path[1:], "/") // omit the first '/' in the Path
	userID := strings.Join(urlParts[1:], "")         // the first element would be 'users'
	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "Bad userID", http.StatusBadRequest)
		return
	}

	colln := senv.DB.Collection("users")

	var resultUser User
	err = colln.FindOne(context.TODO(), bson.D{{Key: "_id", Value: userObjectID}}).Decode(&resultUser)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(writer, "{}")
			return
		}
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	resultUser.PwdHash = "" // set this to empty so that it is not marshalled
	jsonPost, err := json.Marshal(resultUser)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	addCommonHeaders(&writer)
	fmt.Fprintf(writer, string(jsonPost))
}

func (senv *ServerEnv) handlePostCreate(writer http.ResponseWriter, req *http.Request) {
	//fmt.Fprintf(writer, "Hello from handlePostCreate! Invoked at %s.", req.URL.Path[:])

	var post Post

	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	addCommonHeaders(&writer)
	fmt.Fprintf(writer, "Creating post: %v", post)
}

func (senv *ServerEnv) handlePostGet(writer http.ResponseWriter, req *http.Request) {
	urlParts := strings.Split(req.URL.Path[1:], "/")
	postID := strings.Join(urlParts[1:], "")
	postObjectID, err := primitive.ObjectIDFromHex(postID)

	if err != nil {
		http.Error(writer, "Bad userID", http.StatusBadRequest)
		return
	}

	colln := senv.DB.Collection("posts")

	var post Post
	err = colln.FindOne(context.TODO(), bson.D{{Key: "_id", Value: postObjectID}}).Decode(&post)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			fmt.Fprintf(writer, "{}")
			return
		}
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
	}

	jsonPost, err := json.Marshal(post)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	addCommonHeaders(&writer)
	fmt.Fprintf(writer, string(jsonPost))
}

func (senv *ServerEnv) handleUserPostsGet(writer http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(writer, "Hello from handleUserPostsGet! Invoked at %s.", req.URL.Path[:])
}

func main() {
	uri := os.Getenv("MONGODB_URI")
	dbname := os.Getenv("MONGODB_DBNAME")

	if uri == "" || dbname == "" {
		log.Fatal("You must set the MONGODB_URI and MONGODB_DBNAME environment variables.")
	}

	ctx, cancelCtx := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelCtx()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	defer func() {
		log.Println("Closing connection to MongoDB Atlas database")
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	log.Println("Connected to MongoDB Atlas database.")
	log.Println(fmt.Sprintf("Selecting database: %s", dbname))

	senv := &ServerEnv{DB: client.Database(dbname)}
	mux := http.NewServeMux()

	mux.HandleFunc("/users", makeCheckMethodHandler("POST", senv.handleUserCreate))
	mux.HandleFunc("/users/", makeCheckMethodHandler("GET", senv.handleUserGet))
	mux.HandleFunc("/posts", makeCheckMethodHandler("POST", senv.handlePostCreate))
	mux.HandleFunc("/posts/", makeCheckMethodHandler("GET", senv.handlePostGet))
	mux.HandleFunc("/posts/users/", makeCheckMethodHandler("GET", senv.handleUserPostsGet))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
