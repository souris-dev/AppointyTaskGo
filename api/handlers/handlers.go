package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"appyinsta/api/models"
	"appyinsta/api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// All handlers are defined on this struct so that the DB connection
// can be accessed from all the handlers without making a global variable
// This struct can also include other resources that need to be shared
// by the handlers if needed, in the future

type ServerEnv struct {
	DB *mongo.Database
}

// Handlers

// POST /users
func (senv *ServerEnv) HandleUserCreate(writer http.ResponseWriter, req *http.Request) {
	var user models.User

	if err := json.NewDecoder(req.Body).Decode(&user); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.PwdHash == "" || user.Name == "" {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}

	// hash the password of the user
	user.PwdHash = utils.GetHashed256(user.PwdHash)

	colln := senv.DB.Collection("users")
	// ensure that the ID field is empty
	user.UserID = primitive.NilObjectID
	res, err := colln.InsertOne(context.TODO(), user)

	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.UserID = res.InsertedID.(primitive.ObjectID)

	utils.AddCommonHeaders(&writer)
	fmt.Fprintf(writer, "{\"id\": \"%s\"}", user.UserID.Hex())
}

// GET /users/<userID>
func (senv *ServerEnv) HandleUserGet(writer http.ResponseWriter, req *http.Request) {
	urlParts := strings.Split(req.URL.Path[1:], "/") // omit the first '/' in the Path
	userID := strings.Join(urlParts[1:], "")         // the first element would be 'users'
	userObjectID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		http.Error(writer, "Bad userID", http.StatusBadRequest)
		return
	}

	colln := senv.DB.Collection("users")

	var resultUser models.User
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

	utils.AddCommonHeaders(&writer)
	fmt.Fprintf(writer, string(jsonPost))
}

// POST /posts
func (senv *ServerEnv) HandlePostCreate(writer http.ResponseWriter, req *http.Request) {
	var post models.Post

	if err := json.NewDecoder(req.Body).Decode(&post); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	if post.Caption == "" || post.ImgURL == "" || post.PostedByUID == primitive.NilObjectID {
		http.Error(writer, "Bad Request", http.StatusBadRequest)
	}

	// set the PostedOn field of the post as per server time
	post.PostedOn = time.Now().UTC()

	colln := senv.DB.Collection("posts")
	// ensure that the ID field is empty
	post.PostID = primitive.NilObjectID
	res, err := colln.InsertOne(context.TODO(), post)

	if err != nil {
		log.Println(err.Error())
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	post.PostID = res.InsertedID.(primitive.ObjectID)

	utils.AddCommonHeaders(&writer)
	fmt.Fprintf(writer, "{\"id\": \"%s\"}", post.PostID.Hex())
}

// GET /posts/<postID>
func (senv *ServerEnv) HandlePostGet(writer http.ResponseWriter, req *http.Request) {
	urlParts := strings.Split(req.URL.Path[1:], "/")
	postID := strings.Join(urlParts[1:], "")
	postObjectID, err := primitive.ObjectIDFromHex(postID)

	if err != nil {
		http.Error(writer, "Bad userID", http.StatusBadRequest)
		return
	}

	colln := senv.DB.Collection("posts")

	var post models.Post
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

	utils.AddCommonHeaders(&writer)
	fmt.Fprintf(writer, string(jsonPost))
}

// GET /posts/users/<userId>
// This endpoint implements pagination and sends the posts by a user
// in the order of most recent first.
// If in the body, the first_request param is true, then the first n
// number of sorted posts are returned.
// For subsequent requests, the first_request field is either not present or is false
// and the client supplies the last postId and the timestamp that last post it received.
// We query the database for posts that were posted earlier than this timestamp received.

func (senv *ServerEnv) HandleUserPostsGet(writer http.ResponseWriter, req *http.Request) {
	urlParts := strings.Split(req.URL.Path[1:], "/")
	userID := strings.Join(urlParts[2:], "")
	userObjID, err := primitive.ObjectIDFromHex(userID)

	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	var pagInfo models.PostPaginationInfo

	if err := json.NewDecoder(req.Body).Decode(&pagInfo); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var filter bson.D

	if pagInfo.FirstRequest {
		filter = bson.D{
			{Key: "posted_by", Value: userObjID},
		}
	} else {
		filter = bson.D{
			{Key: "posted_by", Value: userObjID},
			{Key: "posted_on", Value: bson.D{{Key: "$gte", Value: pagInfo.LastPostedOn}}},

			// This condition covers an edge case (very unlikely) when two posts have same timestamp:
			{Key: "_id", Value: bson.D{{Key: "$ne", Value: pagInfo.LastPostID}}},
		}
	}

	descendingSort := bson.D{{Key: "posted_by", Value: -1}}
	descendingOpts := options.Find().SetSort(descendingSort).SetLimit(pagInfo.NumberOfNewPosts)

	colln := senv.DB.Collection("posts")
	descendingCursor, descendingErr := colln.Find(context.TODO(), filter, descendingOpts)

	var posts []models.Post
	if descendingErr = descendingCursor.All(context.TODO(), &posts); descendingErr != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	postsJSON, err := json.Marshal(posts)

	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		log.Println(err.Error())
		return
	}

	utils.AddCommonHeaders(&writer)
	fmt.Fprintf(writer, string(postsJSON))
}
