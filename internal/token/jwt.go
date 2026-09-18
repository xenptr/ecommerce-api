package token

import (
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Generator interface {
	Generate(userID int64) (string, error)
}

type Parser interface {
	Parse(tokenString string) (int64, error)
}

type JWT struct {
	secret []byte
}

func NewJWT(secret []byte) *JWT {
	return &JWT{
		secret: secret,
	}
}

func (j *JWT) Generate(userID int64) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := t.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (j *JWT) Parse(tokenString string) (int64, error) {
	var claims jwt.RegisteredClaims

	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
	)

	keyFunc := func(token *jwt.Token) (any, error) {
		return j.secret, nil
	}

	_, err := parser.ParseWithClaims(tokenString, &claims, keyFunc)
	if err != nil {
		return 0, err
	}

	userId, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
