package tool

import "math/rand"

// RandomStringPure generate a random string with length specified.
func RandomStringPure(length int) (result string) {
	buf := make([]byte, length)
	if _, err := rand.Read(buf); err == nil {
		result = string(buf)
	}
	return
	// source:=rand.NewSource(time.Now().UnixNano())
	// b := make([]byte, length)
	// for i := range b {
	// 	b[i] = charset[source.Int63()%int64(len(charset))]
	// }
	// return string(b)
}
