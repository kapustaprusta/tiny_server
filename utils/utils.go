package utils

import (
	"strconv"
	"strings"

	"tiny_server/defs"
)

func ParseDatabaseError(err string) (string, int64) {
	errorParts := strings.Split(err, ":")
	message := errorParts[1]
	code, _ := strconv.ParseInt(strings.Split(errorParts[0], "Error ")[1], 10, 32)

	return message, code
}

func TranslateError(err int64) defs.ErrorMessage {
	var errorMessage defs.ErrorMessage
	switch err {
	case 1062:
		errorMessage.Code = 10
		errorMessage.HTTPCode = 409
		errorMessage.Comment = "Duplicate entry"
	}

	return errorMessage
}
