package core

// Parser is responsible for parsing the error message that uses go template
// syntax and extracting the arguments from it.
// It introduces a few no-op functions to derive the type of the field.
//
// Functions:
// - string
// - int
// - bool
type Parser interface {
	// Parse parses the error message that uses go template syntax and extracts the
	// error parameters and localizable text in go template format.
	// Localizable text leaves only the parameter without its type.
	//
	// Example:
	// Input: "Error with {{ .Param1 | string }}."
	// Output: "Error with {{ .Param1 }}."
	Parse(txt string) (_ map[string]Argument, _ string, err error)
}
