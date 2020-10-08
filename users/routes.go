package users

import (
	"bnr/server"
)

// GetRoutes . . .
func GetRoutes() []server.Route {
	return []server.Route{server.NewRoute("/users", "GET", getUsers)}
}
