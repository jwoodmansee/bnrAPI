package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Server . . .
type Server struct {
	Router *mux.Router
}

// Route . . .
type Route struct {
	Path    string
	Handler func(w http.ResponseWriter, r *http.Request)
	Method  string
	Recover bool
}

type serverError struct {
	Message string
	Status  int
}

// GetIntParam . . .
func GetIntParam(req *http.Request, paramName string) (param int) {
	var paramStrings []string
	var ok bool
	var err error
	fmt.Println(req.URL.Query().Get(paramName))
	if paramStrings, ok = req.URL.Query()[paramName]; !ok || len(paramStrings) != 1 {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s", paramName), http.StatusBadRequest)
	}
	if param, err = strconv.Atoi(paramStrings[0]); err != nil {
		PanicWithStatus(fmt.Errorf("missing or invalid url parameter: %s\nerror: %s", paramName, err.Error()), http.StatusBadRequest)
	}
	return param
}

// PanicWithStatus . . .
func PanicWithStatus(err error, status int) {
	panic(serverError{Message: err.Error(), Status: status})
}

// NewRoute . . .
func NewRoute(path, method string, handler func(w http.ResponseWriter, r *http.Request)) Route {
	return Route{Path: path, Handler: handler, Method: method}
}

func routes() []Route {
	return []Route{}
}

// AddRoutes . . .
func (s *Server) AddRoutes(routes []Route) {
	for _, r := range routes {
		h := s.Router.HandleFunc(r.Path, r.Handler)
		h.Methods(r.Method)
	}
}

// NewServer . . .
func NewServer() *Server {
	s := Server{Router: mux.NewRouter()}
	s.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					switch t := err.(type) {
					case string:
						http.Error(w, t, http.StatusInternalServerError)
					case error:
						http.Error(w, t.Error(), http.StatusInternalServerError)
					case serverError:
						http.Error(w, t.Message, t.Status)
					default:
						http.Error(w, "unknown error", http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		})
	})
	s.AddRoutes(routes())
	return &s
}
