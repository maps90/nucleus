package util

// DefaultString checks if data is nil, returns value instead
func DefaultString(data string, value string) string {
	if data == "" {
		return value
	}

	return data
}

// DefaultInt checks if data is nil, returns value instead
func DefaultInt(data int, value int) int {
	if data == 0 {
		return value
	}

	return data
}

// DefaultFloat checks if data is nil, returns value instead
func DefaultFloat(data float64, value float64) float64 {
	if data == 0 {
		return value
	}

	return data
}
