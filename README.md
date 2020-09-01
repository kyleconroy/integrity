# stakmachine/integrity
[![PkgGoDev](https://pkg.go.dev/badge/github.com/stackmachine/integrity)](https://pkg.go.dev/github.com/stackmachine/integrity) ![GithubActions](https://github.com/stackmachine/integrity/workflows/ci/badge.svg?branch=master)

stackmachine/integrity makes it easy to enable subresource integrity for your web
applications.

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/stackmachine/integrity"
)

func main() {
	// Calculate SHA512 digests for all your static assets
	fs, err := integrity.ParseFiles("static")
	if err != nil {
		log.Fatal(err)
	}

	// Return the digest for a given file path, returning an error if it
	// doesn't exist.
	sha, err := fs.Digest("css/style.css")
	if err != nil {
		log.Fatal(err)
	}

	// Use the parameter on a script or link element
	fmt.Println("integrity", sha)
}
```

The package also ships with a `http.Handler` that checks if an included digest
is valid.

```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stackmachine/integrity"
)

func main() {
	fs, err := integrity.ParseFiles("static")
	if err != nil {
		log.Fatal(err)
	}

	handler := http.FileServer(http.Dir("testdata"))
	handler = integrity.Verify(fs, handler)
	handler = http.StripPrefix("/static/", handler)

	// 200 - GET /static/css/style.css
	// 200 - GET /static/css/style.css?sha=sha512-valid
	// 404 - GET /static/css/style.css?sha=sha512-invalid
	fmt.Println("listening on :8080...")
	http.ListenAndServe(":8080", handler)
}
```
