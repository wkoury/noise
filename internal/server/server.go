package server

import (
	"log"
	"net/http"
)

// net/http handler to serve index.html
func StartHTTPServer() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/server/views/index.html")
	})
	http.HandleFunc("/main.wasm", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "cmd/wasm/main.wasm")
	})
	http.HandleFunc("/wasm_exec.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/server/static/wasm_exec.js")
	})

	// TODO: envvar for port?
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, nil))
}
