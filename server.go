package main

import (
	"encoding/base64"
	"fmt"
	"gconf/github"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Server struct {
	loader         *loader
	snapshot       *configSnapshot
	updateDuration time.Duration
}

func NewServer(repoOwner, repoName string, updateDuration time.Duration) *Server {
	return &Server{
		updateDuration: updateDuration,
		loader: &loader{
			githubAPIHost: github.DefaultHost,
			repoOwner:     repoOwner,
			repoName:      repoName,
			repoSHA:       "master",
		},
	}
}

func (s *Server) ListenAndServe(address string) error {
	// Attempt initial load.
	if err := s.update(); err != nil {
		return err
	}

	// Spawn worker to refresh periodically.
	ticker := time.NewTicker(s.updateDuration)
	quit := make(chan struct{})
	defer close(quit)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := s.update(); err != nil {
					log.Printf("failed to update: %v", err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	r := httprouter.New()
	r.GET("/", s.HandleIndex)
	r.GET("/projects", s.HandleProjects)
	r.GET("/file/*filename", s.HandleFile)
	return http.ListenAndServe(address, r)
}

func (s *Server) HandleIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "gconf")
}

func (s *Server) HandleProjects(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	for _, config := range s.snapshot.configs {
		if strings.Contains(config.path, "/") {
			continue
		}
		fmt.Fprintf(w, "%s\n", config.path)
	}
}

func (s *Server) HandleFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	filename := ps.ByName("filename")
	filename = strings.TrimSpace(filename)
	filename = strings.TrimLeft(filename, "/")
	if filename == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for _, config := range s.snapshot.configs {
		if config.path == filename {
			content, err := base64.StdEncoding.DecodeString(config.blob.Content)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
			}
			maxAge := fmt.Sprintf("max-age=%s", s.updateDuration)
			w.Header().Add("Cache-Control", maxAge)
			w.Header().Add("ETag", config.blob.SHA)
			if _, err := fmt.Fprintf(w, string(content)); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err)
			}
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
}

func (s *Server) update() error {
	snapshot, err := s.loader.Load()
	if err != nil {
		return err
	}
	s.snapshot = snapshot
	return nil
}

type loader struct {
	githubAPIHost string
	repoOwner     string
	repoName      string
	repoSHA       string
}

func (l *loader) Load() (*configSnapshot, error) {
	client := &http.Client{}
	res, err := github.GetTree(client, l.githubAPIHost, l.repoOwner, l.repoName, l.repoSHA, true)
	if err != nil {
		return nil, err
	}

	cs := &configSnapshot{}
	for _, node := range res.Tree {
		blob, err := github.GetBlob(client, l.githubAPIHost, l.repoOwner, l.repoName, node.SHA)
		if err != nil {
			return nil, err
		}
		cs.configs = append(cs.configs, config{
			path: node.Path,
			blob: blob,
		})
	}
	return cs, nil
}

type configSnapshot struct {
	configs []config
}

type config struct {
	path string
	blob *github.GetBlobResponse
}
