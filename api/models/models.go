package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// the struct tags (other than the first one) below are not really required
// but kept here to ease connection between the frontend and backend
// in case the fields in JSON have a different name

type User struct {
	// JSON unmarshalling should skip this field
	// UserID is an objectID provided by mongodb on insertion
	UserID primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string             `json:"name" bson:"name"`
	Email  string             `json:"email" bson:"email"`

	// The following field initially contains the original password just after
	// JSON unmarshalling of the request body, after which it is hashed and updated for storage.
	// Since we do not have a control over the frontend, we're assuming
	// that we are getting the password as plaintext.
	PwdHash string `json:"password,omitempty" bson:"p_hash"`
}

type Post struct {
	PostID      primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	PostedByUID primitive.ObjectID `json:"posted_by" bson:"posted_by"`
	Caption     string             `json:"caption" bson:"caption"`
	ImgURL      string             `json:"img_url" bson:"img_url"`
	PostedOn    time.Time          `json:"posted_on,omitempty" bson:"posted_on"` // filled at the server
}

type PostPaginationInfo struct {
	LastPostID       primitive.ObjectID `json:"last_id"`
	LastPostedOn     time.Time          `json:"last_posted_on"`
	NumberOfNewPosts int64              `json:"n_new"`
	FirstRequest     bool               `json:"first_request,omitempty"`
}
