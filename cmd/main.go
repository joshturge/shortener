package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/joshturge/url-short/pkg/repo"
	"github.com/joshturge/url-short/pkg/server"
)

var (
	listenAddress string
	redisAddress  string
	redisPassword string
	timeout       int
	testRun       bool

	err error
)

func init() {
	flag.StringVar(&listenAddress, "listen", "localhost:8080", "the address the server will listen on")
	flag.BoolVar(&testRun, "test", false, `whether to run the server in test mode which uses an 
	in-memory repository for storing urls`)
	flag.StringVar(&redisAddress, "redis", "localhost:6379", "the address of the redis instance")
	flag.StringVar(&redisPassword, "password", "", "the password of the redis instance")
	flag.IntVar(&timeout, "timeout", 5, "timeout of keys in minutes")
}

func main() {
	flag.Parse()

	// if we are running in test mode then we should use the in-memory keystore
	var rep repo.Repository
	if testRun {
		fmt.Println("INFO: running in test mode")

		rep = repo.NewMap()
		rep.SetTimeout(time.Duration(timeout) * time.Minute)
	} else {
		fmt.Printf("INFO: connecting to the redis instance %s\n", redisAddress)

		rep, err = repo.NewRedisRepo(redisAddress, redisPassword)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			os.Exit(1)
		}
		rep.SetTimeout(time.Duration(timeout) * time.Minute)
	}

	// start the server and wait for a interrupt signal
	srv := server.NewServer(listenAddress, rep)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	fmt.Printf("INFO: serving on %s\n", listenAddress)

	if err := srv.Serve(sigChan); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}

	fmt.Println("INFO: closing connections")

	if err := srv.Close(); err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
	}
}
