package utils

import "regexp"

var objectIDRegex = regexp.MustCompile(`^[a-f0-9]{24}$`)
var technicalNameRegex = regexp.MustCompile(`^[a-z][a-zA-Z0-9]{0,31}$`)

func IsObjectID(id string) bool {
	return objectIDRegex.MatchString(id)
}

func IsTechnicalName(technicalName string) bool {
	return technicalNameRegex.MatchString(technicalName)
}
