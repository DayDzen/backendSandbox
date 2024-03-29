package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	c := getCode("someemail@gmail.com")
	fmt.Println(c)
	c = getCode("somemail@gmail.com")
	fmt.Println(c)
}

func getCode(s string) string {
	h := hmac.New(sha256.New, []byte("ourcode"))
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}
