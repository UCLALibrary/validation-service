package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

const Port = 8888

// App represents your application.
type App struct {
	Router *echo.Echo
}

func main() {
	app := NewApp()
	app.Routes()
	err := app.Router.Start(":" + strconv.Itoa(Port))

	if err != nil {
		panic(err)
	}
}

// NewApp initializes a new instance of your application.
func NewApp() *App {
	e := echo.New()
	return &App{Router: e}
}

// Routes configure the server's endpoints.
func (app *App) Routes() {
	app.Router.GET("/", helloWorld)
	app.Router.RouteNotFound("/*", NotFoundHandler)
}

func helloWorld(c echo.Context) error {
	return c.String(http.StatusOK, "hello world")
}

func NotFoundHandler(c echo.Context) error {
	return echo.NewHTTPError(http.StatusNotFound, "This is not yet supported")
}
