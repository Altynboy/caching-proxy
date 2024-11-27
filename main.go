package main

import (
	"caching-proxy/cache"
	"caching-proxy/components"
	"caching-proxy/response"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

var (
	port   string
	origin string
	clear  bool
	db     *cache.Cache
)

func main() {
	flagParse()
	var err error
	db, err = cache.New("./cache/cache.sql")
	if err != nil {
		log.Fatalf("Error while open db %s", err)
	}
	defer db.Close()

	if clear {
		err := db.DeleteAll()
		if err != nil {
			log.Fatalf("Can't clear cache")
		}
		return
	}

	components.NewHttpClient()
	http.HandleFunc("/", ProxyHandler)

	log.Printf("Server starting on port %s", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	baseUrl := origin
	fmt.Println(r.URL.Path)
	targetUrl, err := url.Parse(baseUrl + r.URL.Path)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	value, ok := db.Get(targetUrl.String())
	if ok {
		err = response.FromCache(w, r, []byte(value))
		if err != nil {
			log.Fatalf("Got err when reading cache %s", err)
			return
		}
	}

	w.Header().Add("X-Cache", "MISS")

	err = components.ProxyClient(targetUrl.String(), w, r, db)
	if err != nil {
		log.Fatalf("Error while proxing %s", err)
	}
}

func flagParse() {
	portFlag := flag.Int("port", 3000, "service port")
	originFlag := flag.String("origin", "https://dummyjson.com", "origin url")
	clearFlag := flag.Bool("clear-cache", false, "clear cache")

	flag.Parse()
	if *portFlag < 1 {
		log.Fatal("Port flag cant be negative")
	}
	port = ":" + strconv.Itoa(*portFlag)
	origin = *originFlag
	clear = *clearFlag
}
