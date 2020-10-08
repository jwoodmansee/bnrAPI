package posts

import (
	"bnr/data"
	"bnr/server"
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

func deletePost(w http.ResponseWriter, r *http.Request) {
	type request struct {
		ID string
	}

	handler := func(req request) {
		db, err := data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		defer db.Disconnect(context.TODO())

		col := db.Database(os.Getenv("Database")).Collection("posts")

		res, err := col.DeleteOne(context.TODO(), bson.A{req.ID})
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		fmt.Println(res)
	}

	var req request = request{ID: mux.Vars(r)["id"]}
	handler(req)
}
