package util

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// ZapLoggerMiddleware integrates the Zap logging into Echo and logs all requests
func ZapLoggerMiddleware(aLogger *zap.Logger) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogLatency:   true,
		LogRemoteIP:  true,
		LogUserAgent: true,
		LogValuesFunc: func(_ echo.Context, values middleware.RequestLoggerValues) error {
			aLogger.Debug("Request",
				zap.String("method", values.Method),
				zap.String("uri", values.URI),
				zap.Int("status", values.Status),
				zap.String("remote_ip", values.RemoteIP),
				zap.String("user_agent", values.UserAgent),
				zap.Duration("latency", values.Latency),
			)
			return nil
		},
	})
}
