package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func main() {
	port := 8080

	http.HandleFunc("/image", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("%s /image route requested", r.Method)

		keys, ok := r.URL.Query()["size"]
		if !ok || len(keys[0]) < 1 {
			http.Error(w, "Missing `size` query string parameter", http.StatusBadRequest)
			return
		}

		sizeInt, err := strconv.Atoi(keys[0])
		if err != nil {
			fmt.Println(err)
			fmt.Fprintln(w, err)
			http.Error(w, "Invalid `size` parameter", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Image Requested - Size %d px", sizeInt)
	})

	fmt.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
