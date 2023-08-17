package client

import (
	"regexp"
)

var (
	roleNotFoundRegExp *regexp.Regexp
	userNotFoundRegExp *regexp.Regexp
)

func init() {
	roleNotFoundRegExp = regexp.MustCompile("auth: Role .* does not exist.")
	userNotFoundRegExp = regexp.MustCompile("auth: User .* does not exist.")
}

// IsKeyNotFound returns true if the error code is ErrorCodeKeyNotFound.
func IsKeyNotFound(err error) bool {
	if cErr, ok := err.(Error); ok {
		return cErr.Code == ErrorCodeKeyNotFound
	}
	return false
}

// IsRoleNotFound returns true if the error means role not found of v2 API.
func IsRoleNotFound(err error) bool {
	if ae, ok := err.(authError); ok {
		return roleNotFoundRegExp.MatchString(ae.Message)
	}
	return false
}

// IsUserNotFound returns true if the error means user not found of v2 API.
func IsUserNotFound(err error) bool {
	if ae, ok := err.(authError); ok {
		return userNotFoundRegExp.MatchString(ae.Message)
	}
	return false
}
