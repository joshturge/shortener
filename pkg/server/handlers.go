package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/joshturge/shortener/pkg/model"
	"github.com/joshturge/shortener/pkg/repo"
	"github.com/joshturge/shortener/pkg/urlhash"
)

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Shorten a url by hashing it and writing it to a repository
func Shorten(rep repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// start checking the request
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if len(r.UserAgent()) < 5 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var user model.Url
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := r.Body.Close(); err != nil {
			fmt.Printf("ERROR: unable to close body of request from client: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if !isUrl(user.URL) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// end checking

		hash, err := urlhash.Hash(user.URL)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err := rep.Set(r.Context(), hash, user.URL); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")

		rHash := model.Hash{Hash: hash}
		// respond with the hash of the url
		if err := json.NewEncoder(w).Encode(&rHash); err != nil {
			fmt.Printf("ERROR: unable to write to response: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

// Index will serve the index of the web application
func Index(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" {
		http.ServeFile(w, r, "./web/static/index.html")
		return
	}
	if r.RequestURI == "/favicon.ico" {
		http.ServeFile(w, r, "./web/static/favicon.ico")
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

// Url will look up the hash of a url in a repository and redirect the client if found
func Url(rep repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// grab the matching strings out of the context of the request
		if hash, ok := r.Context().Value("uri_params").([]string); ok {

			// the first match is most likely the one we want
			url, err := rep.Get(r.Context(), hash[0][1:])
			if err != nil {
				if err == repo.ErrNotFound {
					w.WriteHeader(http.StatusNotFound)
					return
				}

				fmt.Printf("ERROR: %s\n", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			// everything looks good so redirect!
			http.Redirect(w, r, url, http.StatusPermanentRedirect)
			return
		}

		fmt.Printf("ERROR: unable to obtain hash from context\n")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
