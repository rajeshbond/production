package utils

func SafeString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func SafeInt(i *int64) int64 {
	if i == nil {
		return 0
	}

	return *i
}
