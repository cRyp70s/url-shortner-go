package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var sqlitedb = NewDB("urlmaps.db")
var urlmaps = buildInitialMap()

var pageSource = "<html>" +
	"<body><h1>Add URL to shorten</h1><p>%s</p>" +
	"<form action=\"/add\" method=\"POST\">" +
	"<label>From: <input name=\"from\"></label>" +
	"<button>submit</button>" +
	"</form></body>" +
	"</html>"

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, pageSource, "")
}

func getNewURL() string {
	hasher := sha1.New()
	feed := []byte{1, 2, 5, 7, 8, 10}
	for i := range feed {
		rand.Seed(time.Now().UnixNano())
		feed[i] = byte(rand.Int31())
	}
	hasher.Write(feed)
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}

// Insert map into sqlite3 database
func addMap(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		from := strings.Join(r.Form["from"], "")
		to := getNewURL()
		to = to[len(to)-8 : len(to)-2]
		// If user provides url like www.example.com
		// rewrite to https://www.example.com
		if from[:4] == "www." {
			from = "https://" + from
		}
		log.Println("Inserting " + from)
		_, err = sqlitedb.Execute("INSERT INTO urlmaps(original_url, short_url) values (?, ?)", from, to)
		if err != nil {
			fmt.Fprintf(w, "An error occurred: %s", err)
			return
		}
		urlmaps[to] = from
		fmt.Fprintf(w, pageSource, "Your new relative url is: <a href=\"/"+to+"\">/"+to+"</a>")
	} else if r.Method == "GET" {
		fmt.Fprintf(w, pageSource, "")
	}
}

func fallbackHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("/add", addMap)
	return mux
}

func buildInitialMap() map[string]string {
	urlmaps := make(map[string]string)
	rows, _ := sqlitedb.Query("SELECT * FROM urlmaps")
	for rows.Next() {
		var id int
		var original_url string
		var short_url string
		_ = rows.Scan(&id, &original_url, &short_url)
		log.Printf("Redirect found %s -> %s", short_url, original_url)
		urlmaps[short_url] = original_url
	}
	return urlmaps
}

func createHandler() http.HandlerFunc {
	fallback := fallbackHandler()
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request recieved for: %s", r.URL.RequestURI())

		path, ok := urlmaps[r.URL.Path[1:]]
		// Append query params to redirect url
		if len(r.URL.RawQuery) > 0 {
			path = path + "?" + r.URL.RawQuery
		}
		if ok {
			log.Printf("Redirect found %s -> %s", r.URL.Path, path)
			http.Redirect(w, r, path, http.StatusFound)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}
