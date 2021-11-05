package utils

import (
	"github.com/red-gold/telar-core/config"
	"github.com/red-gold/telar-core/types"
	"net/http"
)

func AddPolicies(req *http.Request) http.Request {

	var securityHeaders = map[string]string{
		types.HeaderContentSecurityPolicy: *config.AppConfig.ContentSecurityPolicy,
		types.HeaderReferrerPolicy:        *config.AppConfig.ReferrerPolicy,
		types.HeaderContentTypeOptions:    *config.AppConfig.ContentTypeOptions,
		types.HeaderXSSProtection: 			"X-XSS-Protection: 1; mode=block",
		types.HeaderXFrameOption: 			"SAMEORIGIN",
		types.HeaderHSTS: 			 		"max-age=31536000; includeSubDomains",

	}

	for k, v := range securityHeaders {
		(*req).Header.Add(k, v)
	}

	return *req
}

: 