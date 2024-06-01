package s3ObjectNameGenerator

import "github.com/google/uuid"

func NewObjectName(prefix, extension string) string {
	prefixRunes := []rune(prefix)
	delimiter1 := ""
	delimiter2 := ""
	if prefix != "" && prefixRunes[len(prefixRunes)-1] != '/' {
		delimiter1 = "/"
	}
	if extension != "" && []rune(extension)[0] != '.' {
		delimiter2 = "."
	}
	return prefix + delimiter1 + uuid.New().String() + delimiter2 + extension
}
