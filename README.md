# itsy

itsy is a lightweight web framework written in Go. It is designed to help developers build RESTful APIs that fully embrace the HATEOAS (Hypermedia as the Engine of Application State) principle. itsy provides a simple and intuitive interface for defining resources, handling requests, and rendering responses.

## Why itsy?

RESTful APIs are a popular choice for many web applications due to their simplicity, scalability, and statelessness. However, many APIs do not fully embrace the principles of REST. In particular, they often ignore the HATEOAS principle, which states that a client should be able to navigate an API entirely through hypermedia provided by the server.

itsy aims to make it easy to build truly RESTful APIs by providing built-in support for HATEOAS. With itsy, you can define resources that include links, forms, embedded resources, templates, and actions, allowing clients to navigate your API without any out-of-band information.

## Example

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
  itsy.AddResource(http.MethodGet, "/hello", &HelloWorldResource{})

  // Start the server.
  itsy.Run(":8080")
}

```

In this example, we define a simple resource that returns a "Hello, world!" message when accessed with a GET request. We then add this resource to our itsy application and start the server.
