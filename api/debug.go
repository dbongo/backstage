package api

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

type DebugHandler struct {
	ApiHandler
}

func (handler *DebugHandler) HelloWorld(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse {
	c.Env["Content-Type"] = "text/plain"
	response := &HTTPResponse{StatusCode: http.StatusOK, Payload: "Hello World"}
	return response
}
