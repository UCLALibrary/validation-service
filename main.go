package main

import (
	"errors"
	"fmt"
	"github.com/UCLALibrary/validation-service/api"
	"github.com/UCLALibrary/validation-service/csvutils"
	"github.com/UCLALibrary/validation-service/validation"
	"github.com/UCLALibrary/validation-service/validation/config"
	"github.com/labstack/echo/v4"
	middleware "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
	"log"
	"net/http"
)

// Port is the default port for our server
const Port = 8888

// Service implements the generated API validation interface
type Service struct {
	Engine *validation.Engine
}

// GetStatus handles the GET /status request
func (service *Service) GetStatus(context echo.Context) error {
	// A placeholder response
	return context.JSON(http.StatusOK, api.Status{
		Service: "ok",
		S3:      "ok",
		FS:      "ok",
	})
}

// GetJobStatus handles the GET /status/{jobID} request
func (service *Service) GetJobStatus(context echo.Context, jobID string) error {
	// A placeholder response; structure will change
	return context.JSON(http.StatusOK, api.Status{
		Service: "completed",
		S3:      "ok",
		FS:      "ok",
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
	csvData, readErr := csvutils.ReadUpload(file, logger)
	if readErr != nil {
		return context.JSON(http.StatusBadRequest, map[string]string{"error": "Uploaded CSV file could not be parsed"})
	}

	if err := engine.Validate(profile, csvData); err != nil {
		// Handle if there was a validation error
		return context.JSON(http.StatusCreated, api.Status{
			Service: fmt.Sprintf("error: %v", err),
			S3:      "ok",
			FS:      "ok",
		})
	}

	// Handle if there were no validation errors
	return context.JSON(http.StatusCreated, api.Status{
		Service: "created",
		S3:      "ok",
		FS:      "ok",
	})
}

// Main function starts our Echo server
func main() {
	// Create a new validation engine
	engine, err := validation.NewEngine()
	if err != nil {
		log.Fatal(err)
	}

	// Create the validation service
	service := &Service{
		Engine: engine,
	}

	// Create a new validation application and configure its logger
	echoApp := echo.New()
	echoApp.Use(config.ZapLoggerMiddleware(engine.GetLogger()))

	// Hide the Echo startup messages
	echoApp.HideBanner = true
	echoApp.HidePort = true

	// Handle requests with and without a trailing slash
	echoApp.Pre(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(char echo.Context) error {
			path := char.Request().URL.Path

			if path != "/" && path[len(path)-1] == '/' {
				// Remove trailing slash as our canonical form
				char.Request().URL.Path = path[:len(path)-1]
			}

			return next(char)
		}
	})

	// Load the OpenAPI spec for request validation
	swagger, swaggerErr := api.GetSwagger()
	if swaggerErr != nil {
		log.Fatalf("Failed to load OpenAPI spec: %v", swaggerErr)
	}

	// Register the Echo/OpenAPI validation middleware
	echoApp.Use(middleware.OapiRequestValidator(swagger))

	// Register request handlers
	api.RegisterHandlers(echoApp, service)

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
