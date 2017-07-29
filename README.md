# stakmachine/sri
[![GoDoc](https://godoc.org/stackmachine.com/sri?status.svg)](https://godoc.org/stackmachine.com/sri) [![Build Status](https://travis-ci.org/stackmachine/sri.svg?branch=master)](https://travis-ci.org/stackmachine/sri)

stackmachine/sri makes it easy to enable subresource integrity for your web
applications.

## Install

```
dep ensure stackmachine.com/sri
```

## Usage

```go
package main

import (
    "fmt"

    "stackmachine.com/sri"
)

func main() {
    // Calculate SHA512 digests for all your static assets
    fs, err := sri.ParseFiles("static")
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

The `sri` package also ships with a `http.Handler` that checks if an included
digest is valid.

```go
package main

import (
    "net/http"
    
    "stackmachine.com/sri"
)

fucn main() {
    fs, err := sri.ParseFiles("static")
    if err != nil {
        panic(err)
    }

    handler := http.FileServer(http.Dir("testdata"))
    handler = sri.Verify(handler)
    handler = http.StripPrefix("/static/", handler)
    http.ListenAndServe(handler, nil)
}
```
