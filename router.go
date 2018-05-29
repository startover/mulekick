package mulekick

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// IDEA: Generate swagger API docs based on the calls

type Router struct {
	*mux.Router

	EnableLogging bool
	completeRoute string
	groupRoute    string
	parent        *Router
	middleware    []http.HandlerFunc
}

func (r Router) Get(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("GET", "HEAD")
}

func (r Router) Post(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("POST")
}

func (r Router) Put(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("PUT")
}

func (r Router) Delete(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("DELETE")
}

func (r Router) Patch(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("PATCH")
}

func (r Router) Options(endpoint string, middleware ...http.HandlerFunc) *mux.Route {
	return r.Handle(endpoint, middleware...).Methods("OPTIONS")
}

func (r Router) Handle(endpoint string, mw ...http.HandlerFunc) *mux.Route {
	middleware := make([]http.HandlerFunc, len(r.middleware))
	for i, m := range r.middleware {
		middleware[i] = m
	}

	middleware = append(middleware, mw...)

	if r.parent != nil && len(endpoint) == 0 {
		middleware := middleware[len(r.parent.middleware):]
		return r.parent.Handle(r.groupRoute, middleware...)
	}

	if os.Getenv("ENV") == "debug" {
		fmt.Printf("%v%v (%d handlers)\n", r.completeRoute, endpoint, len(middleware))
	}

	route := r.HandleFunc(endpoint, func(w http.ResponseWriter, req *http.Request) {
		wr := &ResponseWriter{w, false, http.StatusOK, time.Now()}

		for _, m := range middleware {
			m(wr, req)

			if wr.responseWritten {
				break
			}
		}

		if r.EnableLogging {
			wr.LogMiddleware(req)
		}
	})

	return route
}

// Group creates a new sub-router, enabling you to group handlers
func (r Router) Group(str string, middleware ...http.HandlerFunc) Router {
	middleware = append(r.middleware, middleware...)

	newRouter := New(r.PathPrefix(str).Subrouter(), middleware...)
	newRouter.completeRoute = r.completeRoute + str
	newRouter.groupRoute = str
	newRouter.parent = &r
	newRouter.EnableLogging = r.EnableLogging
	return newRouter
}

// Use function adds middleware to the router for calls
func (r *Router) Use(middleware ...http.HandlerFunc) {
	r.middleware = append(r.middleware, middleware...)
}
