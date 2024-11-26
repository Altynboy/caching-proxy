package main

import (
	"caching-proxy/components"
	"log"
	"net/http"
	"net/url"
)

func main() {
	components.NewHttpClient()
	port := ":3000"
	http.HandleFunc("/", ProxyHandler)

	log.Printf("Server starting on port %s", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func ProxyHandler(w http.ResponseWriter, r *http.Request) {
	baseUrl := "https://dummyjson.com"
	targetUrl, err := url.Parse(baseUrl + r.URL.Path)
	if err != nil {
		http.Error(w, "Error parsing URL", http.StatusInternalServerError)
		return
	}

	err = components.ProxyClient(baseUrl+targetUrl.String(), w, r)
	if err != nil {
		log.Fatalf("Error while proxing %s", err)
	}
}
