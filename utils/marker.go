package utils

var characters = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var length = len(characters)

// GetMarkers will generate markers for every list worker.
func GetMarkers(n int) []byte {
	// We only need to split into 62 part as most.
	if n > length {
		n = length
	}
	s := make([]byte, n)

	for i := 0; i < n; i++ {
		s[i] = characters[i*(length/n)]
	}
	return s
}
