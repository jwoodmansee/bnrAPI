package users

import (
	"bnr/data"
	"bnr/server"
	"context"
	"encoding/json"
	"net/http"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func getUsers(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        int    `json:"id" bson:"id"`
		Name      string `json:"name" bson:"name"`
		Email     string `json:"email" bson:"email"`
		Expertise string `json:"expertise" bson:"expertise"`
	}

	handler := func() []byte {
		db, err := data.GetData()
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}

		collection := db.Database(os.Getenv("Database")).Collection("users")
		defer db.Disconnect(context.TODO())
		var res []*response

		cur, err := collection.Find(context.TODO(), bson.D{{}})
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}

		for cur.Next((context.TODO())) {
			var resp response
			err := cur.Decode(&resp)
			if err != nil {
				server.PanicWithStatus(err, http.StatusInternalServerError)
			}

			res = append(res, &resp)
		}
		if err := cur.Err(); err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		cur.Close(context.TODO())
		result, err := json.Marshal(&res)
		if err != nil {
			server.PanicWithStatus(err, http.StatusInternalServerError)
		}
		return result
	}

	w.Header().Set("Content-type", "application/json")
	w.Write(handler())
}
