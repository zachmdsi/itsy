package itsy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	i := New()
	r1 := i.Register("/test/:id")
	r1.GET(func(ctx Context) error {
		id := ctx.GetParam("id")
		ctx.CreateField("testKey", id)
		return ctx.WriteJSON()
	})

	req, err := http.NewRequest("GET", "/test/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	i.ServeHTTP(rr, req)

	expected := `{"fields":{"testKey":"123"},"links":{}}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestWriteHTML(t *testing.T) {
	i := New()
	r1 := i.Register("/test/:id")
	r1.GET(func(ctx Context) error {
		id := ctx.GetParam("id")
		ctx.CreateField("testKey", id)
		return ctx.WriteHTML()
	})

	req, err := http.NewRequest("GET", "/test/123", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	i.ServeHTTP(rr, req)

	expected := `<html>
  <body>
    <div>Fields:
      <div>testKey: 123</div>
    </div>
    <div>Links:
    </div>
  </body>
</html>
`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}
