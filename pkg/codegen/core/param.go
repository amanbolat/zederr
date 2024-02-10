package core

// Param represents a parameter in the error translation message.
//
// Example:
// In message such as "Error with {{ .Param1 | string }}."
// - `Param1` is a name.
// - `string` is a type
type Param struct {
	Name string
	Type string
}
