package jwt

import (
	"crypto/rsa"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

func DecodePrivateKey(privateKey []byte) (*rsa.PrivateKey, error) {
	return jwt.ParseRSAPrivateKeyFromPEM(privateKey)
}

func GenerateToken(duration time.Duration, data map[string]interface{}, key interface{}, issuer string) (string, error) {
	claims := &jwt.MapClaims{
		"iss":  issuer,
		"exp":  time.Now().Add(duration).UTC().Unix(),
		"iat":  time.Now().UTC().Unix(),
		"data": data,
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("RS256"), claims)
	return token.SignedString(key)

}

func ValidateToken(tokenString string, key interface{}) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("RS256") != token.Method {
			return nil, ErrInvalidToken
		}

		return key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
