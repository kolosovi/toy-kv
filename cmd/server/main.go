package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	toykv "github.com/kolosovi/toy-kv"
	"github.com/kolosovi/toy-kv/dao"
	"github.com/kolosovi/toy-kv/internal/walmanager"
)

func Main() error {
	var walFilename string
	flag.StringVar(&walFilename, "wal-filename", "", "")
	flag.Parse()
	if len(walFilename) == 0 {
		return fmt.Errorf("missing argument: wal-filename")
	}
	walManager := walmanager.New(walmanager.WithWALFilename(walFilename))
	d := dao.New(walManager)
	if err := d.Start(); err != nil {
		return fmt.Errorf("dao.Start: %w", err)
	}
	defer func() {
		if err := d.Stop(); err != nil {
			log.Printf("dao.Stop: %v\n", err.Error())
		}
	}()
	http.HandleFunc(
		"/",
		func(w http.ResponseWriter, req *http.Request) {
			key, err := extractKeyFromPath(req.URL.Path)
			if err != nil {
				respondWithError(w, 400, fmt.Errorf("invalid key: %w", err))
				return
			}
			switch req.Method {
			case http.MethodGet:
				value, err := d.Get(toykv.K(key))
				if errors.Is(err, toykv.ErrNotFound) {
					respond(w, 404, fmt.Sprintf("key [%s] not found\n", key))
				} else if err != nil {
					respondWithError(w, 500, fmt.Errorf("dao.Get: %w", err))
				} else {
					respond(w, 200, string(value))
				}
			case http.MethodPut:
				body, err := io.ReadAll(req.Body)
				if err != nil {
					respondWithError(w, 500, fmt.Errorf("failed to read request body: %w", err))
					return
				}
				record := toykv.Record{K: toykv.K(key), V: toykv.V(body)}
				if err := d.Put(record); err != nil {
					respondWithError(w, 500, fmt.Errorf("dao.Put: %w", err))
				} else {
					respond(w, 200, "OK\n")
				}
			case http.MethodDelete:
				err = d.Delete(toykv.K(key))
				if errors.Is(err, toykv.ErrNotFound) {
					respond(w, 404, fmt.Sprintf("key [%s] not found\n", key))
				} else if err != nil {
					respondWithError(w, 500, fmt.Errorf("dao.Delete: %w", err))
				} else {
					respond(w, 200, "OK\n")
				}
			default:
				respondWithNotAllowed(w)
			}
		},
	)
	return http.ListenAndServe(":8181", nil)
}

func extractKeyFromPath(path string) (string, error) {
	if len(path) == 0 {
		return "", fmt.Errorf("empty path")
	}
	path = path[1:]
	if strings.Contains(path, "/") {
		return "", fmt.Errorf("key must not contain slashes")
	}
	if strings.Contains(path, "[") || strings.Contains(path, "]") {
		return "", fmt.Errorf("key must not contain square brackets")
	}
	return path, nil
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	respond(w, code, fmt.Sprintf("%s\n", err))
}

func respond(w http.ResponseWriter, code int, body string) {
	w.WriteHeader(code)
	_, writeErr := w.Write([]byte(body))
	if writeErr != nil {
		log.Printf("error writing response: %v\n", writeErr)
	}
}

func respondWithNotAllowed(w http.ResponseWriter) {
	w.Header().Add("Allow", "GET, PUT, DELETE")
	w.WriteHeader(405)
	w.Write(nil)
}

func main() {
	if err := Main(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
