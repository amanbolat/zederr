package net

import (
	"unicode/utf8"

	"golang.org/x/net/idna"
)

var idnaProfile = idna.New(
	idna.ValidateLabels(true),
	idna.VerifyDNSLength(true),
	idna.StrictDomainName(true),
)

// FQDN reports whether is a valid fully qualified domain name.
//
// FQDN allows only ASCII characters as prescribed by RFC 1034 (A-Z, a-z, 0-9
// and the hyphen).
func FQDN(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] >= utf8.RuneSelf {
			return false
		}
	}

	_, err := idnaProfile.ToASCII(str)

	return err == nil
}
