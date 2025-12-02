package main

import (
    "log"
    "net/http"
)

func main() {
    // Define the file server to serve slides.html
    fileServer := http.FileServer(http.Dir("."))

    // Create a handler for the root path
    http.Handle("/", fileServer)

    // Specify the port to listen on
    port := ":4041"

    // Print a message to indicate the server is running
    log.Printf("Serving slides.html on http://localhost%s", port)

    // Start the server
    err := http.ListenAndServe(port, nil)
    if err != nil {
        log.Fatal("Server error: ", err)
    }
}
