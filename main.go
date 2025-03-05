package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/UCLALibrary/validation-service/api"
	"github.com/UCLALibrary/validation-service/validation"
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/csv"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Port is the default port for our server
const Port = 8888

// ServiceError provides a generic error to use in HTTP responses.
type ServiceError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RouteMapping is a pairing of router path and file system path that can be used to configure request handlers.
type RouteMapping struct {
	RoutePath string
	FilePath  string
}

// TemplateRenderer holds parsed HTML templates for the validation service's Web pages
type TemplateRenderer struct {
	templates *template.Template
	mu        sync.Mutex
}

// Render function on our TemplateRenderer implements Echo's `Renderer` interface
func (renderer *TemplateRenderer) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	renderer.mu.Lock()
	defer renderer.mu.Unlock()
	return renderer.templates.ExecuteTemplate(writer, name, data)
}

// Service implements the generated OpenAPI interface (i.e., handles incoming requests)
type Service struct {
	Engine *validation.Engine
}

// GetStatus handles the GET /status request
func (service *Service) GetStatus(context echo.Context) error {
	// A placeholder response
	return context.JSON(http.StatusOK, api.Status{
		Service:    "ok",
		Fester:     "ok",
		FileSystem: "ok",
	})
}

// UploadCSV handles the /upload/csv POST request
func (service *Service) UploadCSV(context echo.Context) error {
	engine := service.Engine
	logger := engine.GetLogger()

	// Get the CSV file upload and profile
	profile := context.FormValue("profile")
	file, fileErr := context.FormFile("csvFile")
	if fileErr != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": "A CSV file must be uploaded"})
	}

	logger.Debug("Received uploaded CSV file",
		zap.String("csvFile", file.Filename),
		zap.String("profile", profile))

	// Parse the CSV data
	csvData, readErr := csv.ReadUpload(file, logger)

	if readErr != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": "Uploaded CSV file could not be parsed"})
	}

	if err := engine.Validate(profile, csvData); err != nil {
		report, reportErr := csv.NewReport(err, csvData, logger)
		if reportErr != nil {
			logger.Error("Failed to generate report", zap.Error(reportErr), zap.Stack("stacktrace"))

			return context.JSON(http.StatusInternalServerError,
				ServiceError{Code: http.StatusInternalServerError, Message: reportErr.Error()})
		}

		return context.JSON(http.StatusCreated, report)
	}

	return context.JSON(http.StatusCreated, &csv.Report{Profile: profile, Time: time.Now(), Warnings: []csv.Warning{}})
}

// Main function starts our Echo server
func main() {
	// Create a new validation engine for our service to use
	engine, err := validation.NewEngine()
	if err != nil {
		log.Fatal(err)
	}

	// Get the validation engine's logger to use to configure Echo
	logger := engine.GetLogger()

	// Create a new validation application and configure its logger
	echoApp := echo.New()
	echoApp.Use(config.ZapLoggerMiddleware(logger))

	// Hide application startup messages that don't play nicely with logger
	echoApp.HideBanner = true
	echoApp.HidePort = true

	// Turn on Echo's debugging features if we're set to debug (mostly more info in errors)
	if debugging := logger.Check(zap.DebugLevel, "Enable debugging"); debugging != nil {
		echoApp.Debug = true
	}

	// Handle requests with and without a trailing slash using the trailingSlashMiddleware
	echoApp.Pre(trailingSlashMiddleware)

	// Configure the application's route handling
	routes := append(configStaticRoutes(echoApp), configTemplateRoutes(echoApp, getTemplateRenderer(logger))...)
	echoApp.Use(routerConfigMiddleware(echoApp, engine, routes))

	// Log the configured routes when we're running in debug mode
	if debugging := logger.Check(zap.DebugLevel, "Loading routes"); debugging != nil {
		var fields []zap.Field

		for _, route := range echoApp.Routes() {
			routeInfo := []string{route.Method, route.Path}
			fields = append(fields, zap.Strings("route", routeInfo))
		}

		logger.Debug("Registered routes", fields...)
	}

	// Configure the validation server with the port number and the Echo application
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: echoApp,
	}

	// Start the validation server
	if err := echoApp.StartServer(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}

// getTemplateRenderer loads the HTML templates and then provides a template registry that can render them.
func getTemplateRenderer(logger *zap.Logger) *TemplateRenderer {
	templates, templateErr := loadTemplates(logger)
	if templateErr != nil {
		logger.Error(templateErr.Error())
	}

	// Log all the HTML templates that were loaded if we're in debugging mode
	if debugging := logger.Check(zap.DebugLevel, "Load templates"); debugging != nil && templates != nil {
		var fields []zap.Field

		for _, tmpl := range templates.Templates() {
			// In loadTemplates we add a 'root' template with an empty name, but we delete it here it's not used
			if tmpl.Name() != "" {
				fields = append(fields, zap.String("template", tmpl.Name()))
			}
		}

		logger.Debug("Loaded Web resources", fields...)
	}

	return &TemplateRenderer{templates: templates}
}

// loadTemplates loads the available HTML templates for the Web UI.
func loadTemplates(logger *zap.Logger) (*template.Template, error) {
	templates := template.New("") // New set of templates

	// We try both locations: the Docker container's and the local dev's
	patterns := []string{
		// The templates should exist in only one of these locations
		"/usr/local/data/html/template/*.tmpl", // Docker templates
		"html/template/*.tmpl",                 // Local templates
	}

	foundTemplates := false

	// Parse each template path pattern so we can add any matches to the template set
	for _, pattern := range patterns {
		// Check if any files match the pattern before parsing them
		matches, err := filepath.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error checking pattern %s: %w", pattern, err)
		}

		// Skip this pattern if no template files are found here
		if len(matches) == 0 {
			logger.Debug("No templates found (Skipping)", zap.String("location", pattern))
			continue // Moving on to the next location
		}

		// If we get this far, at least one location had templates
		foundTemplates = true

		// Attempt to parse the templates once we know they exist
		templates, err = templates.ParseGlob(pattern)
		if err != nil {
			return nil, fmt.Errorf("error loading templates from %s: %w", pattern, err)
		}
	}

	// If no templates were found in either location, return an error
	if !foundTemplates {
		return nil, fmt.Errorf("no templates found in any of the specified locations")
	}

	// Return the set of template matches
	return templates, nil
}

// configStaticRoutes configures our static resources with the Echo application.
func configStaticRoutes(echoApp *echo.Echo) []RouteMapping {
	staticRoutes := []RouteMapping{
		{"/openapi.yml", "html/assets/openapi.yml"},
		{"/validation.css", "html/assets/validation.css"},
		{"/validation.js", "html/assets/validation.js"},
	}

	for _, route := range staticRoutes {
		echoApp.GET(route.RoutePath, func(aContext echo.Context) error {
			return aContext.File(route.FilePath)
		})
	}

	return staticRoutes
}

// configTemplateRoutes configures our template resources with the Echo application.
func configTemplateRoutes(echoApp *echo.Echo, renderer *TemplateRenderer) []RouteMapping {
	templateRoutes := []RouteMapping{
		{"/", ""},
		{"index.html", ""},
	}

	// Set the Echo application's default template renderer
	echoApp.Renderer = renderer

	// Have the templates renderer handle incoming index requests
	echoApp.GET(templateRoutes[0].RoutePath, func(context echo.Context) error {
		data := map[string]interface{}{
			"Version": os.Getenv("VERSION"),
		}

		return context.Render(http.StatusOK, templateRoutes[1].RoutePath, data)
	})

	return templateRoutes
}

// routerConfigMiddleware configures the application's router with a fully configured OpenAPI set of routes.
func routerConfigMiddleware(echoApp *echo.Echo, engine *validation.Engine, routes []RouteMapping) echo.MiddlewareFunc {
	swagger, swaggerErr := api.GetSwagger()
	if swaggerErr != nil {
		engine.GetLogger().Fatal("Failed to load OpenAPI spec", zap.Error(swaggerErr))
	}

	// Register OpenAPI defined request handlers for our service
	api.RegisterHandlers(echoApp, &Service{
		Engine: engine,
	})

	// We return the oapi-codegen middleware that handles our OpenAPI defined routes
	return middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Skipper: func(aContext echo.Context) bool {
			for index := range routes {
				// We ignore paths that we've already configured through static or template handlers
				if aContext.Path() == routes[index].RoutePath {
					return true
				}
			}

			return false
		},
	})
}

// trailingSlashMiddleware handles paths with slashes at the end so they also resolve.
func trailingSlashMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		path := context.Request().URL.Path

		// Strip trailing slashes if found in path
		if path != "/" && path[len(path)-1] == '/' {
			context.Request().URL.Path = path[:len(path)-1]
		}

		return next(context)
	}
}
