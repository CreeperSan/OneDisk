package formatstring

import (
	"OneDisk/lib/random"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func String(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func Password(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(hashedPassword)
}

func GenerateRefreshToken() string {
	return random.String(32)
}

func GenerateToken() string {
	return random.String(32)
}
