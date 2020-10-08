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

func getUserPosts(w http.ResponseWriter, r *http.Request) {
	type postsResponse struct {
		ID     int    `json:"id" bson:"id"`
		UserID int    `json:"userId" bson:"userId"`
		Title  string `json:"title" bson:"title"`
		Body   string `json:"body" bson:"body"`
	}
	type response struct {
		Posts []postsResponse `json:"posts"`
	}

	type request struct {
		UserID int
	}

	handler := func(req request) *response {
		var err error
		var db *mongo.Client
		db, err = data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		defer db.Disconnect((context.TODO()))

		collection := db.Database(os.Getenv("Database")).Collection("posts")
		cur, err := collection.Find(context.TODO(), bson.D{{"userId", req.UserID}})
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		defer cur.Close(context.TODO())

		var posts response
		for cur.Next((context.TODO())) {
			var post postsResponse
			err := cur.Decode(&post)
			if err != nil {
				server.PanicWithStatus(err, http.StatusInternalServerError)
			}
			posts.Posts = append(posts.Posts, post)
		}
		return &posts
	}

	params := strings.Split(r.URL.Path, "/")
	userID, err := strconv.Atoi(params[2])
	if err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}
	var req request = request{UserID: userID}

	w.Header().Set("Content-type", "application/json")
	if err = json.NewEncoder(w).Encode(handler(req)); err != nil {
		server.PanicWithStatus(err, http.StatusInternalServerError)
	}
}
