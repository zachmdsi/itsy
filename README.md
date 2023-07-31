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

type HelloWorldResource struct {
  itsy.BaseResource
}

func (h *HelloWorldResource) Render() string {
  return "<h1>Hello, world!</h1>"
}

func main() {
  // Create a new itsy instance.
  itsy := itsy.New()

  // Add a resource.
  itsy.Register("/hello", &HelloWorldResource{})

  // Start the server.
  itsy.Run(":8080")
}

```

In this example, we define a simple resource that returns a "Hello, world!" message when accessed. We then add this resource to our itsy application and start the server.

## Contributing

We welcome contributions to itsy! If you'd like to contribute, please open an issue or submit a pull request.
