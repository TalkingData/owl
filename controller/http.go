package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func InitHttpServer() (err error) {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok = r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}
		}
	}()
	go startHttpServer()
	return nil
}

func startHttpServer() error {
	router := NewRouter()
	return http.ListenAndServe(GlobalConfig.HTTP_SERVER, router)
}

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandleFunc
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}

	return router
}

type Route struct {
	Name       string
	Method     string
	Pattern    string
	HandleFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
}

func Index(w http.ResponseWriter, r *http.Request) {
	for _, node := range controller.nodePool.Nodes {
		fmt.Fprintln(w, string(node.Encode())+",")
	}
}
