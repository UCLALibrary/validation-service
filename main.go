package main

import (
	"errors"
	"fmt"
	"github.com/UCLALibrary/validation-service/api"
	"github.com/UCLALibrary/validation-service/validation"
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/UCLALibrary/validation-service/validation/utils"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// Port is the default port for our server
const Port = 8888

// Service implements the generated API validation interface
type Service struct {
	Engine *validation.Engine
}

// TemplateRegistry holds parsed HTML templates for the validation service's Web pages
type TemplateRegistry struct {
	templates *template.Template
	mu        sync.Mutex
}

// Render implements Echo's `Renderer` interface
func (tmplReg *TemplateRegistry) Render(writer io.Writer, name string, data interface{}, context echo.Context) error {
	tmplReg.mu.Lock()
	defer tmplReg.mu.Unlock()
	return tmplReg.templates.ExecuteTemplate(writer, name, data)
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
	csvData, readErr := utils.ReadUpload(file, logger)
	if readErr != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": "Uploaded CSV file could not be parsed"})
	}

	if err := engine.Validate(profile, csvData); err != nil {
		// Handle if there was a validation error
		return context.JSON(http.StatusCreated, api.Status{
			Service:    fmt.Sprintf("error: %v", err),
			Fester:     "ok",
			FileSystem: "ok",
		})
	}

	// Handle if there were no validation errors
	return context.JSON(http.StatusCreated, api.Status{
		Service:    "created",
		Fester:     "ok",
		FileSystem: "ok",
	})
}

// Main function starts our Echo server
func main() {
	// Create a new validation engine for our service to use
	engine, err := validation.NewEngine()
	if err != nil {
		log.Fatal(err)
	}

	// Logger we can use to output information
	logger := engine.GetLogger()

	// Create a new validation application and configure its logger
	echoApp := echo.New()
	echoApp.Use(config.ZapLoggerMiddleware(engine.GetLogger()))

	// Hide Echo startup messages that don't play nicely with logger
	echoApp.HideBanner = true
	echoApp.HidePort = true

	// Turn on Echo's debugging features if we're set to debug (mostly more info in errors)
	if debugging := logger.Check(zap.DebugLevel, "Enable debugging"); debugging != nil {
		echoApp.Debug = true
	}

	// Serves the service's OpenAPI specification at the expected endpoint
	echoApp.GET("/openapi.yml", func(context echo.Context) error {
		return context.File("html/assets/openapi.yml")
	})

	// Sets the template renderer for the application
	echoApp.Renderer = setTemplateRenderer(logger)

	// Handle our SPA endpoint outside of the OpenAPI specification
	echoApp.GET("/", func(context echo.Context) error {
		data := map[string]interface{}{
			"Version": "0.0.1",
		}

		return context.Render(http.StatusOK, "index.html", data)
	})

	// Handle requests with and without a trailing slash
	echoApp.Pre(TrailingSlashMiddleware)

	// Load the OpenAPI spec for request validation
	swagger, swaggerErr := api.GetSwagger()
	if swaggerErr != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", swaggerErr)
	}

	// Register the Echo/OpenAPI validation middleware; have it ignore things served independent of the OpenAPI spec
	echoApp.Use(middleware.OapiRequestValidatorWithOptions(swagger, &middleware.Options{
		Skipper: func(context echo.Context) bool {
			return context.Path() == "/openapi.yml" || context.Path() == "/"
		},
	}))

	// Register request handlers for our service
	api.RegisterHandlers(echoApp, &Service{
		Engine: engine,
	})

	// Log the configured routes when we're running in debug mode
	if debugging := logger.Check(zap.DebugLevel, "Loading routes"); debugging != nil {
		var fields []zap.Field

		for _, route := range echoApp.Routes() {
			routeInfo := []string{route.Method, route.Path}
			fields = append(fields, zap.Strings("route", routeInfo))
		}

		logger.Debug("Registered routes", fields...)
	}

	// Configure the validation echoApp
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", Port),
		Handler: echoApp,
	}

	// Start the validation echoApp
	if err := echoApp.StartServer(server); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}

// TrailingSlashMiddleware handles paths with slashes at the end so they also resolve.
func TrailingSlashMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(context echo.Context) error {
		path := context.Request().URL.Path

		// Strip trailing slashes if found in path
		if path != "/" && path[len(path)-1] == '/' {
			context.Request().URL.Path = path[:len(path)-1]
		}

		return next(context)
	}
}

// SetTemplateRenderer loads the HTML templates and then provides a template registry that can render them.
func setTemplateRenderer(logger *zap.Logger) *TemplateRegistry {
	templates, templateErr := loadTemplates(logger)
	if templateErr != nil {
		logger.Error(templateErr.Error())
	}

	// Log all the HTML templates that were loaded if we're in debugging mode
	if debugging := logger.Check(zap.DebugLevel, "Load templates"); debugging != nil && templates != nil {
		var fields []zap.Field

		for _, tmpl := range templates.Templates() {
			// We add a 'root' logger with an empty name, but it does nothing so delete it
			if tmpl.Name() != "" {
				fields = append(fields, zap.String("template", tmpl.Name()))
			}
		}

		logger.Debug("Loaded Web resources", fields...)
	}

	return &TemplateRegistry{templates: templates}
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
