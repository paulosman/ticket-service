package main

import (
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
)

func main() {
    router := mux.NewRouter()
    router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]bool{"ok": true})
    })

    server := &http.Server{
        Addr: "0.0.0.0:9000",
        Handler: router,
        ReadTimeout: 10 * time.Second,
    }

    log.Fatal(server.ListenAndServe())
}
