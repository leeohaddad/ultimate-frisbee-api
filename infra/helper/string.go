package helper

import "fmt"

func IsNilOrEmpty(str *string) bool {
	return str == nil || *str == ""
}

func ErrorMessageInField(entity, field string) string {
	return fmt.Sprintf("the %s's '%s' should not be empty", entity, field)
}
