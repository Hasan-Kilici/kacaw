# HTTP Methods
All HTTP Methods in Kacaw
::code-group
```go [main.go]
package main

import (
	"fmt"
	"github.com/Hasan-Kilici/kacaw"
	"net/http"
)

func main() {
	// Create a new instance of the Kacaw router
	router := kacaw.Default()

	// Handle CONNECT request to "/connect" URL
	router.CONNECT("/connect", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a CONNECT request")
	})

	// Handle DELETE request to "/users/{id}" URL
	router.DELETE("/users/:id", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a DELETE request")
	})

	// Handle GET request to the root URL ("/")
	router.GET("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a GET request")
	})

	// Handle HEAD request to "/info" URL
	router.HEAD("/info", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a HEAD request")
	})

	// Handle OPTIONS request to "/options" URL
	router.OPTIONS("/options", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is an OPTIONS request")
	})

	// Handle POST request to "/data" URL
	router.POST("/data", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a POST request")
	})

	// Handle PUT request to "/users/{id}" URL
	router.PUT("/users/:id", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a PUT request")
	})

	// Handle TRACE request to "/trace" URL
	router.TRACE("/trace", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "This is a TRACE request")
	})

	// Run the server on port 8000
	router.Run(":8000")
}
```
::
