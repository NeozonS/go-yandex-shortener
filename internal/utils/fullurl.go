package utils

import (
	"strings"
)

func FullURL(baseurl, token string) string {
	if !strings.HasSuffix(baseurl, "/") {
		baseurl += "/"
	}
	return baseurl + token
}
