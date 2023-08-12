# itsy

itsy is a lightweight web framework written in Go, designed to facilitate the creation of hypermedia-driven web applications. itsy is built with the principles of hypermedia navigation at its core, providing a robust foundation for building web applications that fully leverage the power of hypermedia.

## Why itsy?

Web APIs have become a cornerstone of modern web development, with RESTful APIs being particularly popular due to their simplicity, scalability, and statelessness. However, many APIs fall short of fully embracing the principles of REST, often overlooking the HATEOAS principle. This principle posits that a client should be able to navigate an API entirely through hypermedia provided by the server, without needing any out-of-band information.

itsy is designed to address this gap, making it straightforward to build web applications that harness the full potential of hypermedia. With itsy, you can define resources that include links, forms, embedded resources, templates, and actions, enabling clients to navigate your API seamlessly and intuitively.

## Simple Example

Here's a simple example of how to use itsy to create a RESTful API:

```go
package main

import (
  "net/http"
  "itsy"
)

func main() {
  i := itsy.New()

  r1 := i.Register("/main/:id")
  r2 := i.Register("/linked/:id")

  r1.GET(func(ctx itsy.Context) error {
    id := ctx.GetParam("id")
    ctx.CreateField("Main", id)
    return ctx.WriteHTML()
  })

  r2.GET(func(ctx itsy.Context) error {
    id := ctx.GetParam("id")
    ctx.CreateField("Linked", id)
    return ctx.WriteHTML()
  })

  r1.Link(r2, "related")

  i.Run(":8080")
}

```

In this example, we create two resources, `/main/:id` and `/linked/:id`. The first resource is the main resource, and the second is a related resource. We then define a handler for the main resource that writes a string to the response. Finally, we link the main resource to the related resource, and run the server.

When we run the server and navigate to `http://localhost:8080/main/1`, we see the following response:

```html
HTTP/1.1 200 OK

<html>
  <body>
    <div>Fields:
      <div>Main: 1</div>
    </div>
    <div>Links:
      <a href="/linked/1">related</a>
    </div>
</html>
```

## Contributing

**Itsy is still in early development.**

We welcome contributions to itsy! If you'd like to contribute, please open an issue or submit a pull request.
