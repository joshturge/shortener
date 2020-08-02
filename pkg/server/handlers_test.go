package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshturge/url-short/pkg/repo"
	"github.com/joshturge/url-short/pkg/server"
)

var hashTest = map[string]string{
	"03e400d9": "https://joshturge.dev",
	"cd919b1e": "https://google.com",
	"a9ecbd6c": "https://drive.google.com",
}

var rep repo.Repository

func init() {
	rep = repo.NewMap()
	for hash, url := range hashTest {
		rep.Set(context.Background(), hash, url)
	}
}

func TestShorten(t *testing.T) {
	for hash, data := range hashTest {
		buf := bytes.Buffer{}
		if err := json.NewEncoder(&buf).
			Encode(struct {
				Url string `json:"url"`
			}{data}); err != nil {
			t.Errorf("unable to encode test user data: %w", err)
			t.FailNow()
		}

		req := httptest.NewRequest("POST", "http://localhost:8080/shorten", &buf)
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("User-Agent", "some-browser")

		w := httptest.NewRecorder()
		handler := server.Shorten(rep)

		handler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("response was not 200 got: %s", w.Result().Status)
			t.FailNow()
		}

		var rHash struct {
			Hash string `json:"hash"`
		}

		if err := json.NewDecoder(w.Body).Decode(&rHash); err != nil {
			t.Errorf("failed to encode test response data into json: %w", err)
			t.FailNow()
		}

		if rHash.Hash != hash {
			t.Errorf("hash does not match wanted: %s got: %s", hash, rHash.Hash)
		}
	}
}
