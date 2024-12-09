package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vitaminmoo/poe-get-version/internal/version"
	"google.golang.org/appengine"
	"google.golang.org/appengine/memcache"
)

const (
	cacheFor = 60 * time.Second
)

func main() {
	http.HandleFunc("/", indexHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "GET")
	if r.Method == http.MethodOptions {
		return
	}

	if r.Method == http.MethodGet {
		versions, cached, err := getCache(ctx, "poe-versions", getVersions)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", int(cacheFor.Seconds())))
		w.Header().Add("X-From-Cache", fmt.Sprint(cached))
		fmt.Fprintln(w, versions)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func getVersions() (string, error) {
	poe, err := version.Poe()
	if err != nil {
		log.Printf("getting poe version: %v", err)
		poe = "error"
	}
	poe2, err := version.Poe2()
	if err != nil {
		log.Printf("getting poe2 version: %v", err)
		poe2 = "error"
	}
	// marshal poe and poe2 into json
	output := map[string]string{
		"poe":  poe,
		"poe2": poe2,
	}

	json, err := json.Marshal(output)
	if err != nil {
		fmt.Println(err)
	}
	return string(json), nil
}

func setCache(ctx context.Context, key string, value string) error {
	item := &memcache.Item{
		Key:   key,
		Value: []byte(value),
	}
	return memcache.Set(ctx, item)
}

func getCache(ctx context.Context, key string, fn func() (string, error)) (string, bool, error) {
	item, err := memcache.Get(ctx, key)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			value, err := fn()
			if err != nil {
				return "", false, fmt.Errorf("calling fn: %w", err)
			}
			item := &memcache.Item{
				Key:        key,
				Value:      []byte(value),
				Expiration: cacheFor,
			}
			err = memcache.Set(ctx, item)
			if err != nil {
				log.Printf("setting cache: %v", err)
			}
			return value, false, nil
		}
		return "", false, err
	}
	return string(item.Value), true, nil
}
