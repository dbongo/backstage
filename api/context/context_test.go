package context

import (
  "net/http"
  "net/http/httptest"
  "testing"
  "encoding/json"

  . "gopkg.in/check.v1"
  "github.com/zenazn/goji/web"
  "github.com/albertoleal/backstage/errors"
)

func Test(t *testing.T) { TestingT(t) }

type S struct {}

var _ = Suite(&S{})

func (s *S) TestAddGetRequestError(c *C) {
  m := web.New()

  m.Get("/helloworld", func(c web.C, w http.ResponseWriter, r *http.Request) {
    AddRequestError(&c, &errors.HTTPError{StatusCode: http.StatusUnauthorized,
      Message: "You do not have access to this resource."})

    key, _ := GetRequestError(&c)
    body, _ := json.Marshal(key)
    http.Error(w, string(body), key.StatusCode)
  })

  req, _ := http.NewRequest("GET", "/helloworld", nil)
  recorder := httptest.NewRecorder()
  env := map[string]interface{}{}
  m.ServeHTTPC(web.C{Env: env}, recorder, req)

  c.Assert(recorder.Code, Equals, 401)
  c.Assert(recorder.Body.String(), Equals, "{\"status_code\":401,\"message\":\"You do not have access to this resource.\",\"url\":\"\"}\n")
}