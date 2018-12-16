# stakmachine/integrity
[![GoDoc](https://godoc.org/stackmachine.com/integrity?status.svg)](https://godoc.org/stackmachine.com/integrity) [![Build Status](https://travis-ci.org/stackmachine/integrity.svg?branch=master)](https://travis-ci.org/stackmachine/integrity)

stackmachine/integrity makes it easy to enable subresource integrity for your web
applications.

## Install

```
dep ensure github.com/stackmachine/integrity
```

## Usage

```go
package main

import (
    "fmt"

    "github.com/stackmachine/integrity"
)

func main() {
    // Calculate SHA512 digests for all your static assets
    fs, err := integrity.ParseFiles("static")
    if err != nil {
        panic(err)
    }

    // Return the digest for a given file path, returning an error if it
    // doesn't exist.
    sha, err := fs.Digest("css/style.css")
    if err != nil {
        panic(err)
    }

    // Set the `integrity` parameter on a script or link element
    fmt.Printf(`<script type="javascript" integrity="%s" src="...">`, sha)
}
```

The `integrity` package also ships with a `http.Handler` that checks if an included
digest is valid.

```go
package main

import (
    "net/http"
    
    "github.com/stackmachine/integrity"
)

func main() {
    fs, err := integrity.ParseFiles("static")
    if err != nil {
        panic(err)
    }

    handler := http.FileServer(http.Dir("testdata"))
    handler = integrity.Verify(handler, fs)
    handler = http.StripPrefix("/static/", handler)

    // 200 - GET /static/css/style.css 
    // 200 - GET /static/css/style.css?sha=sha512-valid
    // 404 - GET /static/css/style.css?sha=sha512-invalid
    http.ListenAndServe(handler, nil)
}
```
