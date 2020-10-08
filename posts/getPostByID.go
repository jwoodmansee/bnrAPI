package posts

import (
	"bnr/data"
	"bnr/server"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func getPostByID(w http.ResponseWriter, r *http.Request) {
	type postsResponse struct {
		ID     int    `json:"id" bson:"id"`
		UserID int    `json:"userId" bson:"userId"`
		Title  string `json:"title" bson:"title"`
		Body   string `json:"body" bson:"body"`
	}

	type request struct {
		ID int
	}

	handler := func(req request) *postsResponse {
		var err error
		var db *mongo.Client
		db, err = data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		defer db.Disconnect((context.TODO()))

		collection := db.Database(os.Getenv("Database")).Collection("posts")
		var post postsResponse
		err = collection.FindOne(context.TODO(), bson.D{{"id", req.ID}}).Decode(&post)
		return &post
	}

	var err error
	postID := strings.Split(r.URL.Path, "/")

	id, err := strconv.Atoi(postID[2])
	if err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}

	var req request = request{ID: id}

	w.Header().Set("Content-type", "application/json")
	if err = json.NewEncoder(w).Encode(handler(req)); err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}
}
