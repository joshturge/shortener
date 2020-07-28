package server

// TODO(joshturge): maybe move this into its own project in the future?

import (
	"context"
	"net/http"
	"regexp"
)

// Route defines a handler that gets called when the corresponding regex pattern matches
type Route struct {
	pattern *regexp.Regexp
	handle  http.HandlerFunc
}

// Router will route requests to their respective handlers
type Router struct {
	routes        []Route
	defaultHandle http.HandlerFunc
}

// NewRouter creates a new router with a default handle of 404
func NewRouter() *Router {
	return &Router{defaultHandle: func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}}
}

// SetDefaultHandle will set the default handle to h
func (ro Router) SetDefaultHandle(h http.HandlerFunc) {
	ro.defaultHandle = h
}

// HandlerFunc will register a new route within the router given the
// regex pattern to match and the handle to call
func (ro *Router) HandlerFunc(pattern string, handle http.HandlerFunc) {
	route := Route{pattern: regexp.MustCompile(pattern), handle: handle}

	ro.routes = append(ro.routes, route)
}

// ServeHTTP satisfies the http.Handle interface. It will call a routes handler if its regex
// matches against the requests URL path
func (ro *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range ro.routes {
		if matches := route.pattern.FindStringSubmatch(r.URL.Path); len(matches) > 0 {
			// store all the matching strings in the context of the request
			ctx := context.WithValue(r.Context(), "uri_params", matches)
			r = r.WithContext(ctx)

			route.handle(w, r)
			return
		}
	}

	ro.defaultHandle(w, r)
}
