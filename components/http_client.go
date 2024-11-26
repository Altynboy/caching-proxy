package components

import (
	"io"
	"log"
	"net/http"
)

var httpClient *http.Client

func NewHttpClient() {
	httpClient = &http.Client{}
}

func ProxyClient(url string, w http.ResponseWriter, r *http.Request) error {
	req, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		return err
	}
	req.Header = make(http.Header)
	for h, val := range r.Header {
		req.Header[h] = val
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for h, val := range resp.Header {
		w.Header()[h] = val
	}

	w.WriteHeader(resp.StatusCode)

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Proxied %s %s - Status: %d", r.Method, err, resp.StatusCode)
	return nil
}
