package posts

import (
	"bnr/data"
	"bnr/server"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func createPost(w http.ResponseWriter, r *http.Request) {
	type request struct {
		UserID int    `json:"userId"`
		Title  string `json:"title"`
		Body   string `json:"body"`
	}

	handler := func(req request) {
		db, err := data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		defer db.Disconnect(context.TODO())

		collection := db.Database(os.Getenv("Database")).Collection("posts")

		res, err := collection.InsertOne(context.TODO(), bson.M{"userId": req.UserID, "title": req.Title, "body": req.Body})
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		fmt.Println(res.InsertedID)
	}

	var err error
	var req request

	if err = json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.PanicWithStatus(err, http.StatusBadRequest)
	}

	handler(req)
}
