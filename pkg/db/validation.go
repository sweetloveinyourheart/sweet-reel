package db

import "regexp"

func validatePrefix(prefix string) bool {
	pattern := `\A[[:alpha:]_-]+\z`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(prefix)
}
