# Get Started
Step 1: Open your terminal and execute the command "go get github.com/Hasan-Kilici/kacaw".
::code-group
  ```bash [Terminal]
  go get github.com/Hasan-Kilici/kacaw
  ```
::
Step 2: In your main.go file, add the following code:
::code-group
```go [main.go]
package main

import (
	"fmt"
	"net/http"

	"github.com/Hasan-Kilici/kacaw"
)

func main() {
	r := kacaw.Default()

	r.GET("/", func(w http.ResponseWriter, req *http.Request) {
		data := map[string]interface{}{
			"Title": "Hello, World!",
		}
		r.HTML(http.StatusOK, "index.html", data, w)
	})

	r.Run(":8000")
}
```
```html [index.html]
<html>
    <head>
        <title>{{.Title}}</title>
    </head>
    <body>
        <h1>{{.Title}}</h1>
    </body>
</html>
```
::
The above code sets up a basic server using Kacaw. It listens for incoming requests on the root URL ("/") and responds with an HTML template. In this example, the template displays a title "Hello, World!".

To explain the purpose of this code, it establishes a minimal web server using the Kacaw framework. When a GET request is made to the root URL ("/"), the server renders an HTML template named "index.html" with the provided data (in this case, just the title). The resulting HTML response will display the title "Hello, World!".

Please note that you need to have the necessary HTML template file ("index.html") in the appropriate directory for this code to work properly.

Feel free to adjust and expand upon this code to suit your specific project needs. Happy coding with Kacaw!