package csv

import (
	"encoding/json"
	"errors"
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"strings"
	"time"
)

// Warning is an individual validation warning.
type Warning struct {
	Message  string `json:"message"`
	Header   string `json:"header"`
	ColIndex int    `json:"column"`
	RowIndex int    `json:"row"`
	Value    string `json:"value"`
}

// Report is a collection of validation warnings.
type Report struct {
	Profile  string    `json:"profile"`
	Time     time.Time `json:"time"`
	Warnings []Warning `json:"warnings"`
}

// NewReport creates a report of validation warnings.
func NewReport(multiErr error, csvData [][]string, logger *zap.Logger) (*Report, error) {
	report := &Report{}

	// Set the time the report was generated
	report.Time = time.Now()

	// Cycle through the csv.Error(s) and add them to the report
	for _, csvErr := range multierr.Errors(multiErr) {
		var err *Error

		ok := errors.As(csvErr, &err)
		if ok {
			location := err.Location

			// Set report's profile if it's not already been set
			if report.Profile == "" {
				report.Profile = err.Profile
			}

			header, headerErr := GetHeader(location, csvData, report.Profile)
			if headerErr != nil {
				// At this point in the process, this shouldn't be able to happen
				logger.Error("header error", zap.Error(headerErr), zap.Stack("stacktrace"))
			}

			report.Warnings = append(report.Warnings, Warning{
				strings.ReplaceAll(err.String(), "\n", "<br/>"),
				header,
				err.Location.ColIndex, // The front-end should make this 1-based
				err.Location.RowIndex, // The front-end should make this 1-based
				strings.ReplaceAll(csvData[location.RowIndex][location.ColIndex], "\n", "\\n"),
			})
		} else {
			logger.Error("Unexpected error", zap.Error(err), zap.Stack("stacktrace"))
		}
	}

	return report, nil
}

// SerializeReport serializes the Report to JSON for return to the Web browser.
func SerializeReport(report *Report) (string, error) {
	jsonData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
