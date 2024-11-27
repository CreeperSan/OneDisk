package random

import "math/rand"

func Int(max int) int {
	return rand.Intn(max)
}

func Charset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[Int(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return Charset(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
}

func StringWithSymbol(length int) string {
	return Charset(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
}

func Password(length int) string {
	return Charset(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()")
}
