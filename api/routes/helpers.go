package routes

import "strconv"

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func intPtr(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

func int64Ptr(value int64) *int64 {
	if value == 0 {
		return nil
	}
	return &value
}

func parseIntPtr(value string) *int {
	if value == "" {
		return nil
	}
	num, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &num
}
