package main

import (
	"Codenames/internal/common"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type IRouter interface {
	InitRouter() *http.ServeMux
}

type router struct{}

func (router *router) InitRouter() *http.ServeMux {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", Chain(SampleUsage, Logging()))
	serveMux.HandleFunc("/CreateGameWithParameter", Chain(CreateGameWithParameter,  Method("POST"), Logging()))
	return serveMux
}

func CreateGameWithParameter(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var p common.StartGame
	err := decoder.Decode(&p)
	fmt.Println(p)
	if err != nil {
		panic(err)
	}
	log.Println(p)
	common.Response(p)
}


func SampleUsage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Sample Usage\n/Get?key=XXX\n/Set?key=XXX&value=XXX")
	fmt.Fprintln(w, "/Flush\n/Delete?key=XXX")
	fmt.Fprintln(w, "/GetFrequency\n/SetFrequency?Frequency=XXX")
	fmt.Fprintln(w, "/GetPath\n/SetPath?Path=XXX")
	fmt.Fprintln(w, "/GetImageOfMemory")
}


var (
	myRouter	*router
	routerOnce	sync.Once
)

func Router() IRouter {
	routerOnce.Do(func() {
		myRouter = &router{}
	})
	return myRouter
}

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path and the time it took to process
func Logging() Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware things
			start := time.Now()
			defer func() { log.Println(r.URL.Path, time.Since(start)) }()

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Method ensures that url can only be requested with a specific method, else returns a 400 Bad Request
func Method(m string) Middleware {

	// Create a new Middleware
	return func(f http.HandlerFunc) http.HandlerFunc {

		// Define the http.HandlerFunc
		return func(w http.ResponseWriter, r *http.Request) {

			// Do middleware check
			if r.Method != m {
				keys, _ := r.URL.Query()["key"]
				fmt.Println(keys)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			// Call the next middleware/handler in chain
			f(w, r)
		}
	}
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	for _, m := range middlewares {
		f = m(f)
	}
	return f
}
