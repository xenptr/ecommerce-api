package token

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

const tokenTypeRefresh = "refresh"

type refreshClaims struct {
	jwt.RegisteredClaims
	TokenType string `json:"type"`
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
	return j.sign(baseClaims(userID, accessTokenTTL))
}

func (j *JWT) GenerateRefresh(userID int64) (string, error) {
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := refreshClaims{
		RegisteredClaims: baseClaims(userID, refreshTokenTTL),
		TokenType:        tokenTypeRefresh,
	}

	claims.ID = jti

	return j.sign(claims)
}

func (j *JWT) Parse(tokenString string) (int64, error) {
	var claims jwt.RegisteredClaims

	if err := j.parse(tokenString, &claims); err != nil {
		return 0, err
	}

	return strconv.ParseInt(claims.Subject, 10, 64)
}

func (j *JWT) ParseRefresh(tokenString string) (int64, error) {
	var claims refreshClaims

	if err := j.parse(tokenString, &claims); err != nil {
		return 0, err
	}

	if claims.TokenType != tokenTypeRefresh {
		return 0, errors.New("invalid token type")
	}

	return strconv.ParseInt(claims.Subject, 10, 64)
}

func (j *JWT) sign(claims jwt.Claims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(j.secret)
}

func (j *JWT) parse(tokenString string, claims jwt.Claims) error {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"HS256"}),
	)

	keyFunc := func(token *jwt.Token) (any, error) {
		return j.secret, nil
	}

	_, err := parser.ParseWithClaims(tokenString, claims, keyFunc)
	return err

}

func baseClaims(userID int64, ttl time.Duration) jwt.RegisteredClaims {
	now := time.Now()

	return jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
}

// generateJTI returns a cryptographically random 16-byte hex string.
func generateJTI() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
