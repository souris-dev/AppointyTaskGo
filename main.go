package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"appointyGo/api"
	"appointyGo/api/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	senv := &api.ServerEnv{DB: client.Database(dbname)}
	mux := http.NewServeMux()

	mux.HandleFunc("/users", utils.MakeCheckMethodHandler("POST", senv.HandleUserCreate))
	mux.HandleFunc("/users/", utils.MakeCheckMethodHandler("GET", senv.HandleUserGet))
	mux.HandleFunc("/posts", utils.MakeCheckMethodHandler("POST", senv.HandlePostCreate))
	mux.HandleFunc("/posts/", utils.MakeCheckMethodHandler("GET", senv.HandlePostGet))
	mux.HandleFunc("/posts/users/", utils.MakeCheckMethodHandler("GET", senv.HandleUserPostsGet))

	log.Fatal(http.ListenAndServe(":8080", mux))
}
