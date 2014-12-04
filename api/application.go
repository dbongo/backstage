package api

import (
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/tsuru/config"
	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

type Application struct {
}

func (app *Application) Init() {
	err := config.ReadConfigFile("config.yaml")

	if err != nil {
		fmt.Printf("Error reading config file: %s\n", err.Error())
	}
}

func (app *Application) DrawRoutes() {
	goji.Use(RequestIdMiddleware)
	goji.NotFound(NotFoundHandler)

	// Controllers
	servicesController := &ServicesController{}
	debugController := &DebugController{}
	usersController := &UsersController{}
	groupsController := &GroupsController{}

	// Public Routes
	goji.Get("/", app.Route(servicesController, "Index"))
	goji.Post("/api/users", app.Route(usersController, "CreateUser"))
	goji.Post("/api/signin", app.Route(usersController, "SignIn"))
	goji.Use(ErrorHandlerMiddleware)

	// Private Routes
	api := web.New()
	goji.Handle("/api/*", api)
	api.Use(middleware.SubRouter)
	api.NotFound(NotFoundHandler)
	api.Use(AuthorizationMiddleware)
	api.Get("/helloworld", app.Route(debugController, "HelloWorld"))
	api.Delete("/users", app.Route(usersController, "DeleteUser"))

	api.Post("/teams", app.Route(groupsController, "CreateTeam"))
	api.Delete("/teams/:id", app.Route(groupsController, "DeleteTeam"))
}

func (app *Application) Route(controller interface{}, route string) interface{} {
	fn := func(c web.C, w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		methodValue := reflect.ValueOf(controller).MethodByName(route)
		methodInterface := methodValue.Interface()
		method := methodInterface.(func(c *web.C, w http.ResponseWriter, r *http.Request) *HTTPResponse)
		response := method(&c, w, r)
		w.WriteHeader(response.StatusCode)
		if _, exists := c.Env["Content-Type"]; exists {
			w.Header().Set("Content-Type", c.Env["Content-Type"].(string))
		}
		io.WriteString(w, response.Payload)
	}
	return fn
}