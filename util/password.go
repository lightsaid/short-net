package util

import "golang.org/x/crypto/bcrypt"

func GenHashedPassword(plainText string) (string, error) {
	buf, err := bcrypt.GenerateFromPassword([]byte(plainText), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func CheckHashedPassword(hashedPass string, plainText string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(plainText))
}
