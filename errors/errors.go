// Package errors contains error messages that are used by the validators.
package errors

// Error messages
var (
	NilProfileErr        = "supplied profile cannot be nil"
	NoPrefixErr          = "ARK must start with 'ark:/'"
	NaanTooShortErr      = "NAAN must be at least 5 digits long"
	NaanProfileErr       = "The supplied NAAN is not allowed for the supplied profile"
	NoObjIDErr           = "The ARK must contain an object identifier"
	InvalidObjIDErr      = "The object identifier and qualifier is not valid"
	ArkValFailed         = "ARK validation failed"
	EolFoundErr          = "character for EOL found in cell"
	BadHeaderErr         = "could not retrieve CSV header: %s"
	FieldNotFoundErr     = "required field `%s` was not found"
	FieldDataNotFoundErr = "data for required field `%s` was not found"
	UnknownProfileErr    = "unknown profile `%s`"
	ProfileConfigErr     = "supplied profile has objTypes and notObjTypes set: %s"
	NoHostDir            = "a HOST_DIR must be set"
	FileNotExist         = "the file path given does not exist: %s"
	URLFormatErr         = "license URL is not in a proper format (check for HTTPS)"
	URLConnectErr        = "problem connecting to license URL"
	URLReadErr           = "problem reading body of license URL"
	URLDupeBadErr        = "duplicate invalid license URL"
	TypeWhitespaceError  = "field contains invalid characters (e.g., spaces, line breaks)"
	TypeValueError       = "object type field doesn't contain valid value"
	VisibilityValueError = "visibility field doesn't contain valid value"
	NotAnIntErr          = "the Item Sequence is not an integer"
	NotAPosIntErr        = "the Item Sequence value is not a positive integer"
	UnicodeErr           = "field contains unicode replacement char (�)"
	DupeUnicodeErr       = "field duplicates earlier entry with unicode replacement chari (�)"
)
