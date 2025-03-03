package main

import (
	"net/http"
)

func main() {
	const port = "8080"
	const filePathRoot = "."

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	svr := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	svr.ListenAndServe()
}
