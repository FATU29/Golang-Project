package utils

func StringValue(str *string) string {
	if str == nil {
		return ""
	}
	return *str
}
