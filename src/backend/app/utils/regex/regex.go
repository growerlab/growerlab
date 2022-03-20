package regex

import "regexp"

var PasswordRegex = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
var UsernameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
var RepositoryNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-\.]{2,50}$`)

func Match(val string, reg *regexp.Regexp) bool {
	return reg.MatchString(val)
}
