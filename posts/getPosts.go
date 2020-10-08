package posts

import (
	"bnr/data"
	"bnr/server"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func getAllPosts(w http.ResponseWriter, r *http.Request) {
	type postsResponse struct {
		ID     int    `json:"id" bson:"id"`
		UserID int    `json:"userId" bson:"userId"`
		Title  string `json:"title" bson:"title"`
		Body   string `json:"body" bson:"body"`
	}

	handler := func() []byte {
		db, err := data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}

		collection := db.Database(os.Getenv("Database")).Collection("posts")
		defer db.Disconnect(context.TODO())
		var posts []*postsResponse

		cur, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}

		for cur.Next((context.TODO())) {
			var p postsResponse
			err := cur.Decode(&p)
			if err != nil {
				server.PanicWithStatus(err, http.StatusInternalServerError)
			}
			posts = append(posts, &p)
		}
		if err := cur.Err(); err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		cur.Close(context.TODO())
		p, err := json.Marshal(posts)
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		return p
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(handler())
}
