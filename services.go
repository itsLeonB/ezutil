package ezutil

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/itsLeonB/ezutil/config"
	"github.com/itsLeonB/ezutil/internal"
	"github.com/rotisserie/eris"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	Data map[string]any `json:"data"`
}

type JWTService interface {
	CreateToken(data map[string]any) (string, error)
	VerifyToken(tokenstr string) (JWTClaims, error)
}

type jwtServiceHS256 struct {
	issuer        string
	secretKey     string
	tokenDuration time.Duration
}

func NewJwtService(configs *Auth) JWTService {
	return &jwtServiceHS256{
		issuer:        configs.Issuer,
		secretKey:     configs.SecretKey,
		tokenDuration: configs.TokenDuration,
	}
}

func (j *jwtServiceHS256) CreateToken(data map[string]any) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    j.issuer,
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.tokenDuration)),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
			Data: data,
		},
	)

	signed, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", eris.Wrap(err, "error signing token")
	}

	return signed, nil
}

func (j *jwtServiceHS256) VerifyToken(tokenstr string) (JWTClaims, error) {
	var claims JWTClaims

	_, err := jwt.ParseWithClaims(
		tokenstr,
		&claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(j.secretKey), nil
		},
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, UnauthorizedError(config.MsgAuthExpiredToken)
		}

		return claims, eris.Wrap(err, "error parsing token")
	}

	return claims, nil
}

type HashService interface {
	Hash(val string) (string, error)
	CheckHash(hash, val string) (bool, error)
}

func NewHashService(cost int) HashService {
	if cost < 0 {
		cost = 10 // TODO: make this configurable
	}
	return &internal.HashServiceBcrypt{Cost: cost}
}
