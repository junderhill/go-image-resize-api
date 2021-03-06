package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/nfnt/resize"
)

func main() {
	port := 8080

	r := mux.NewRouter()

	r.HandleFunc("/image/{file}", imageHandler)
	http.Handle("/", r)

	fmt.Printf("Listening on port %d\n", port)

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + strconv.Itoa(port),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	if vars["file"] == "" {
		fmt.Printf("No file specified\n")
		w.WriteHeader(http.StatusNotFound)
	}
	filename := vars["file"]

	fmt.Printf("%s %s requested \n", r.Method, r.URL)

	keys, ok := r.URL.Query()["size"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "Missing `size` query string parameter\n", http.StatusBadRequest)
		return
	}

	sizeInt, err := strconv.Atoi(keys[0])
	if err != nil {
		fmt.Println(err)
		fmt.Fprintln(w, err)
		http.Error(w, "Invalid `size` parameter\n", http.StatusBadRequest)
		return
	}

	resizedImage := getImage("images/"+filename, uint(sizeInt))
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
