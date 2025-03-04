package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/phongnd2802/daily-social/internal/global"
)

type PayloadClaims struct {
	jwt.RegisteredClaims
}

func generateTokenJWT(payload jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(global.Config.SecretKeyJWT))
}

func CreateToken(uuidToken string, expirationTime string) (string, error) {
	expiration, err := time.ParseDuration(expirationTime)
	if err != nil {
		return "", err
	}

	now := time.Now()
	expiresAt := now.Add(expiration)
	return generateTokenJWT(PayloadClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			ExpiresAt: &jwt.NumericDate{Time: expiresAt},
			IssuedAt:  &jwt.NumericDate{Time: now},
			Issuer:    "Daily Social",
			Subject:   uuidToken,
		},
	})
}

func ParseJwtTokenSubject(token string) (*jwt.RegisteredClaims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.RegisteredClaims{}, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(global.Config.SecretKeyJWT), nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*jwt.RegisteredClaims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

func VerifyTokenSubject(token string) (*jwt.RegisteredClaims, error) {
	claims, err := ParseJwtTokenSubject(token)
	if err != nil {
		return nil, err
	}
	if err = claims.Valid(); err != nil {
		return nil, err
	}
	return claims, nil
}
