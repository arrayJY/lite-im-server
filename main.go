package main

import (
	"github.com/arrayJY/go-im-server/services"
	"github.com/emicklei/go-restful/v3"
	"log"
	"net/http"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

func main() {
	restful.Add(services.AuthService())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
