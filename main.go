package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"time"
)

// Command-line flags
var (
	repoOwner, repoName string
	port                int
	updateSecs          int
)

func init() {
	flag.StringVar(&repoOwner, "repo-owner", "", "The name of the  Git repository owner")
	flag.StringVar(&repoName, "repo-name", "", "The name of the Git repository")
	flag.IntVar(&port, "port", 8080, "The port to listen on")
	flag.IntVar(&updateSecs, "t", 0, "The number of seconds to wait between config updates")
}

func checkFlags() error {
	if repoOwner == "" {
		return errors.New("-repo-owner is required")
	}
	if repoName == "" {
		return errors.New("-repo-name is required")
	}
	if port <= 0 {
		return fmt.Errorf("invalid port number: %d", port)
	}
	if updateSecs <= 0 {
		return errors.New("-t must a positive integer")
	}
	return nil
}

func execute() error {
	address := fmt.Sprintf(":%d", port)
	log.Printf("listening at %s", address)
	s := NewServer(repoOwner, repoName, time.Duration(updateSecs)*time.Second)
	return s.ListenAndServe(address)
}

func main() {
	flag.Parse()
	if err := checkFlags(); err != nil {
		log.Fatal(err)
	}
	if err := execute(); err != nil {
		log.Fatal(err)
	}
}
