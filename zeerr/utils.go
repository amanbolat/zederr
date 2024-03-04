package zeerr

import (
	"path"
)

// MakeUID creates a unique error identifier.
// Format: <domain>/<namespace>/<code>
// Example: "acme.com/auth/invalid_credentials"
func MakeUID(domain, namespace, code string) string {
	return path.Join(domain, namespace, code)
}
