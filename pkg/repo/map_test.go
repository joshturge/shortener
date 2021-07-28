package repo_test

import (
	"context"
	"time"
	"testing"

	"github.com/joshturge/shortener/pkg/repo"
)

var hashTest = map[string]string{
	"58b70814": "https://joshturge.dev",
	"403e06b6": "https://google.com",
	"706608fe": "https://drive.google.com",
}

var ctx context.Context

func init() {
	ctx = context.Background()
}

func TestSet(t *testing.T) {

	var (
		rep = repo.NewMap(time.Duration(3) * time.Second)
		err error
	)

	for key, val := range hashTest {
		if err = rep.Set(ctx, key, val); err != nil {
			t.Errorf("unable to set %s to %s: %s", key, val, err.Error())
		}
	}

	var val string

	for key, tVal := range hashTest {
		if val, err = rep.Get(ctx, key); err != nil {
			t.Errorf("unable to get value from key: %s: %s", key, err.Error())
		}

		if val != tVal {
			t.Errorf("unexpected value wanted: %s got: %s", tVal, val)
		}
	}
}


func TestGet(t *testing.T) {

	var (
		rep = repo.NewMap(time.Second)
	)

	TestSet(t)

	_, err := rep.Get(ctx, "doesn't exist")
	if err != repo.ErrNotFound {
		t.Errorf("unexpected error wanted ErrNotFound got: %s", err.Error())
	}

	if err = rep.Set(ctx, "should expire", "val"); err != nil {
			t.Errorf("unable to set key/value pair: %s", err.Error())
	}

	time.Sleep(time.Second)

	_, err = rep.Get(ctx, "should expire")
	if err != repo.ErrNotFound {
		t.Errorf("unexpected error wanted ErrNotFound got: %s", err.Error())
	}
}
