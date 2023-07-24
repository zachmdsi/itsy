# itsy

**itsy** is a minimalistic, hypermedia-driven web framework in Go, focusing on simplicity and ease of use while offering core features such as routing, middleware, and logging. itsy is intended to be lightweight yet highly flexible, supporting the development of robust RESTful APIs with hypermedia controls.

_itsy is under active development and is not fully intended for production use._

## Quickstart

Install the package:

```bash
go get github.com/zachmdsi/itsy
```

## Code Example

```go
package main

import (
  "github.com/zachmdsi/itsy"
)

func main() {
  i := itsy.New()

  i.GET("/books/:title", func(ctx itsy.Context) error {
    // Fetch book data and create a new Book resource
    title := ctx.Request().URL.Query().Get("title")
    book := itsy.NewBook(title, "John Doe")

    // Render the Book resource as JSON
    return ctx.RenderResource(book)
  })

  i.Run(":8080")
}
```

## Contribute

Contributions are welcome. Please submit a pull request or create an issue for any enhancements you think of.
