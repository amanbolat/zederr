package zederr

import (
	"net"
	"strings"
)

#NamespaceName: string & =~"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" & strings.MaxRunes(255) & strings.MinRunes(1)
#ArgumentName: string & =~"^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$" & strings.MaxRunes(255) & strings.MinRunes(1)
#ErrorName: string & =~"^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$" & strings.MaxRunes(255) & strings.MinRunes(1)
#LocaleName: "en" | "zh"

#schema_version: string
#domain: net.FQDN & string
#namespaces: [#NamespaceName]: close({
	description: string
	errors: [#ErrorName]: close({
		description: string
		deprecated: string
		status_code: {
			http: >= 0 & <= 599
			grpc: >= 1 & <= 16
		}
		arguments: [#ArgumentName]: close({
			description!: string
			deprecated?: string
			is_internal?: bool
			type!: "string" | "int" | "float" | "bool" | "timestamp"
			translations: {[#LocaleName]: string}
		})
		message: {
			internal!: string
			public: {[#LocaleName]: string}
		}
		solution?: {
			internal!: string
			public: {[#LocaleName]: string}
		}
	})
})

#Spec: {
	schema_version: #schema_version
	domain: #domain
	namespaces: #namespaces
}
