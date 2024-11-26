package main

import (
	"caching-proxy/components"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// TODO: Cache response and add header
// TODO: Implement clear cache
var (
	port   string
	origin string
)

func main() {
	flagParse()
	components.NewHttpClient()
	http.HandleFunc("/", ProxyHandler)

	log.Printf("Server starting on port %s", port)
	err := http.ListenAndServe(port, nil)
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

	err = components.ProxyClient(targetUrl.String(), w, r)
	if err != nil {
		log.Fatalf("Error while proxing %s", err)
	}
}

func flagParse() {
	portFlag := flag.Int("port", 3000, "service port")
	originFlag := flag.String("origin", "https://dummyjson.com", "origin url")

	flag.Parse()
	if *portFlag < 1 {
		log.Fatal("Port flag cant be negative")
	}
	port = ":" + strconv.Itoa(*portFlag)
	origin = *originFlag
}
