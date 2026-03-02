package services

import "golang.org/x/crypto/bcrypt"

// HashPassword хеширует пароль с помощью bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPassword сравнивает пароль с хешем
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// MustHash хеширует пароль или паникует (для демо-данных)
func MustHash(pwd string) string {
	h, _ := HashPassword(pwd)
	return h
}
