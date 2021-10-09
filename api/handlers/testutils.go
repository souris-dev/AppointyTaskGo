package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func openDBCon() (*ServerEnv, *mongo.Client) {
	uri := os.Getenv("MONGODB_URI")
	dbname := os.Getenv("MONGODB_DBNAME")

	if uri == "" || dbname == "" {
		log.Fatal("You must set the MONGODB_URI and MONGODB_DBNAME environment variables.")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	log.Println("Connected to MongoDB Atlas database.")
	log.Println(fmt.Sprintf("Selecting database: %s", dbname))

	senv := &ServerEnv{DB: client.Database(dbname)}
	return senv, client
}

func checkResponseHeaders(resp *http.Response) error {
	ctype := resp.Header.Get("Content-Type")

	if ctype == "" {
		return fmt.Errorf("Content-Type not set.")
	} else if ctype == "application/json; charset=utf-8" {
		return nil
	} else {
		return fmt.Errorf("Content-Type found was: %s", ctype)
	}
}
