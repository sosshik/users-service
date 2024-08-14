package utils

import "strings"

func ProcessFilter(filter string) (string, string) {
	fieldAndValue := strings.Split(filter, "=")
	if len(fieldAndValue) == 2 {
		return fieldAndValue[0], fieldAndValue[1]
	}
	return "", ""
}
