package response

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
)

func FromCache(w http.ResponseWriter, r *http.Request, cache []byte) error {

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(cache)), r)
	if err != nil {
		return err
	}

	for h, val := range resp.Header {
		for _, val := range val {
			w.Header().Add(h, val)
		}
	}
	w.Header().Set("X-Cache", "HIT")
	w.WriteHeader(resp.StatusCode)

	io.Copy(w, resp.Body)
	resp.Body.Close()
	return nil
}
