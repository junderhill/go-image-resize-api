package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/nfnt/resize"
)

func main() {
	port := 8080

	http.HandleFunc("/image", imageHandler)

	fmt.Printf("Starting server on port %d", port)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
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

	resizedImage := getImage("images/01.jpg", uint(sizeInt))
	fmt.Fprintf(w, "%s", resizedImage)
}

func getImage(filename string, size uint) []byte {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	m := resize.Resize(size, 0, img, resize.Lanczos3)

	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, m, nil); err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
