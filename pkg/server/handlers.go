package server

import (
	"encoding/hex"
	"fmt"
	"hash/adler32"
	"net/http"
	"net/url"

	"github.com/joshturge/shortener/pkg/repo"
)

const uriKeyName = "url"

func isUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

// Shorten a url by hashing it and writing it to a repository
func Shorten(rep repo.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost  {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		url := r.PostFormValue(uriKeyName)

		if !isUrl(url) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// create the url hash
		chksm := adler32.Checksum([]byte(url))
		hash := hex.EncodeToString([]byte{
			byte(chksm >> 24),
			byte(chksm >> 16),
			byte(chksm >> 8),
			byte(chksm),
		})

		var err error

		if err = rep.Set(r.Context(), hash, url); err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "text/plain")

		var n int
		if n, err = w.Write([]byte(hash)); n != len(hash) || err != nil {
			fmt.Printf("ERROR: %s: wrote %d/%d bytes\n", err.Error(), n, len(url))
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
