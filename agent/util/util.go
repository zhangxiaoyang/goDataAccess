package util

import (
	"regexp"
)

func ExtractDomain(url string) string {
	return regexp.MustCompile(`http[s]?://([\w\-\.]+)`).FindStringSubmatch(url)[1]
}
