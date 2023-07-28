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
  itsy.Add(http.MethodGet, "/hello", &HelloWorldResource{})

  // Start the server.
  itsy.Run(":8080")
}

```

In this example, we define a simple resource that returns a "Hello, world!" message when accessed with a GET request. We then add this resource to our itsy application and start the server.

## Advanced Example

```go
package main

import (
  "html/template"
  "net/http"
  "itsy"
)

type PostResource struct {
  itsy.BaseResource
  Title string
  Body  string
}

func (p *PostResource) Render(ctx itsy.Context) string {
  t := template.Must(template.New("post").Parse(`
    <h1>{{.Title}}</h1>
    <p>{{.Body}}</p>
    <a href="/comments?post={{.Title}}">View Comments</a>
  `))

  var rendered string
  err := t.Execute(&rendered, map[string]string{"Title": p.Title, "Body": p.Body})
  if err != nil {
    panic(err)
  }

  return rendered
}

type CommentResource struct {
  itsy.BaseResource
  Author string
  Text   string
}

func (c *CommentResource) Render(ctx itsy.Context) string {
  postTitle := ctx.Params()["post"]

  t := template.Must(template.New("comment").Parse(`
    <h3>Comment by {{.Author}}</h3>
    <p>{{.Text}}</p>
    <a href="/post/{{.PostTitle}}">Back to Post</a>
  `))

  var rendered string
  err := t.Execute(&rendered, map[string]string{"Author": c.Author, "Text": c.Text, "PostTitle": postTitle})
  if err != nil {
    panic(err)
  }

  return rendered
}

func main() {
  // Create a new itsy instance.
  itsy := itsy.New()

  // Add a post resource.
  postResource := &PostResource{Title: "My First Post", Body: "This is my first blog post."}
  itsy.Add(http.MethodGet, "/post/:title", postResource)

  // Add a comment resource.
  commentResource := &CommentResource{Author: "John Doe", Text: "Great post!"}
  commentResource.AddLink(itsy.A("/post/My%20First%20Post", "Back to Post"))
  itsy.Add(http.MethodGet, "/comments", commentResource)

  // Start the server.
  itsy.Run(":8080")
}
```

In this example, a GET request to /post/My%20First%20Post would return the blog post titled "My First Post", along with a link to view its comments. A GET request to /comments?post=My%20First%20Post would return a comment by "John Doe", along with a link back to the post it belongs to.

## Contributing

We welcome contributions to itsy! If you'd like to contribute, please open an issue or submit a pull request.
