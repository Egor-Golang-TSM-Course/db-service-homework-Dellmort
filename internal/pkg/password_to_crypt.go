package pkg

import (
	"crypto/sha1"
	"fmt"
)

func GenetareSHA1(password string) string {
	h := sha1.New()
	h.Write([]byte(password))

	return fmt.Sprintf("%x", h.Sum(nil))
}
