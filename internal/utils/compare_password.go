package utils

import "golang.org/x/crypto/bcrypt"

type ComparePasswordInterface interface {
	ComparePassword(hashedPassword, inputPassword string) error
}

type ComparePass struct{}

func NewPasswordUtil() *ComparePass {
	return &ComparePass{}
}

func (h *ComparePass)  ComparePassword(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}
