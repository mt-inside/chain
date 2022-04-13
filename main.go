package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		spew.Dump(r)
		msg := os.Getenv("CHAIN_OUTPUT")
		if msg == "" {
			log.Fatal("Must supply output to return")
		}
		fmt.Fprintf(w, "%s\n", msg)

		next := os.Getenv("CHAIN_NEXT")
		if next == "" {
			return
		}
		client := &http.Client{}
		req, err := http.NewRequest("GET", "http://"+next+"/", nil)
		if err != nil {
			log.Fatal(err)
		}
		for k, _ := range r.Header {
			if strings.HasPrefix(k, "X-Override-") {
				req.Header.Add(k, r.Header.Get(k))
			}
		}
		spew.Dump(req)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(w, string(body))
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
