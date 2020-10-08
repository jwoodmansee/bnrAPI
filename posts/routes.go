package posts

import (
	"bnr/server"
)

//GetRoutes . . .
func GetRoutes() []server.Route {
	return []server.Route{server.NewRoute("/posts", "GET", getAllPosts),
		server.NewRoute("/post/{postId}", "GET", getPostByID),
		server.NewRoute("/user/{id}/posts", "GET", getUserPosts)}
}
