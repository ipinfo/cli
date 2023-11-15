package lib

import "regexp"

var DomainRegex = regexp.MustCompile(`(?:[a-zA-Z0-9-]+\.){1,}[a-zA-Z]{2,}|(?:[^\s.]+\.)+[^\s]{2,}`)
