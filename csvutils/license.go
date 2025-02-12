package csvutils

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type License struct {
	location Location
	license  string
}

func IsValidLicense(license string) error {
	r, _ := regexp.Compile("^http\\:\\/\\/[0-9a-zA-Z]([-.\\w]*[0-9a-zA-Z])*(:(0-9)*)*(\\/?)([a-zA-Z0-9\\-\\.\\?\\,\\'\\/\\\\+&amp;%\\$#_]*)?$")
	if !r.MatchString(license) {
		return fmt.Errorf("License URL %s in cell [%d][%d] not in a proper format", license, location.RowIndex, location.ColIndex)
	}

	resp, err := http.Get(license)
	if err != nil {
		return fmt.Errorf("Error connecting to license URL in cell [%d][%d]: %s", err.Error(), location.RowIndex, location.ColIndex)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading body of license URL in cell [%d][%d]: %s", err.Error(), location.RowIndex, location.ColIndex)
	}
	if len(body) == 0 {
		return fmt.Errorf("License URL %s  in cell [%d][%d] appears to lack content", license, location.RowIndex, location.ColIndex)
	}

	// Supplied license is valid
	return nil
}
