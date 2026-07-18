package main

import (
	"fmt"
	"io"
	"net/http"
)

// http.ResponseWriter is used to control the response information
// written back to the client that made the request, such as the
// body of the response or the status code.

// Then, the *http.Request value is used to get information about
// the request that came into the server, such as the body being
// sent in the case of a `POST` request of information about the
// the client that made the request.
func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is ze website\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}
