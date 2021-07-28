package server_test

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"net/http/httptest"
	"time"
	"testing"
	"fmt"

	"github.com/joshturge/shortener/pkg/repo"
	"github.com/joshturge/shortener/pkg/server"
)

var hashTest = map[string]string{
	"58b70814": "https://joshturge.dev",
	"403e06b6": "https://google.com",
	"706608fe": "https://drive.google.com",
}

var rep repo.Repository

func init() {
	rep = repo.NewMap(time.Duration(3) * time.Second)
	for hash, url := range hashTest {
		rep.Set(context.Background(), hash, url)
	}
}

func TestShorten(t *testing.T) {
	for hash, rurl := range hashTest {
		buf := bytes.Buffer{}
		if _, err := buf.WriteString(fmt.Sprintf("url=%s", url.PathEscape(rurl))); err != nil {
			t.Error(err.Error())
			t.FailNow()
		}
			
		req := httptest.NewRequest("POST", "http://localhost:8080/shorten", &buf)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

		w := httptest.NewRecorder()
		handler := server.Shorten(rep)

		handler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("response was not 200 got: %s", w.Result().Status)
			t.FailNow()
		}

		/*
		var rHash struct {
			Hash string `json:"hash"`
		}

		if err := json.NewDecoder(w.Body).Decode(&rHash); err != nil {
			t.Errorf("failed to encode test response data into json: %w", err)
			t.FailNow()
		}
		*/

		if bytes.Compare(w.Body.Bytes(), []byte(hash)) != 0 {
			t.Errorf("response body does not match expected hash: %s got %s",
				hash, string(w.Body.Bytes()))
		}
	}
}
