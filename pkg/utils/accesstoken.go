package utils

import (
	"errors"
	"fmt"
	"github.com/Zhoangp/Api_Gateway-Courses/config"
	"github.com/Zhoangp/Api_Gateway-Courses/pkg/common"
	"github.com/golang-jwt/jwt/v4"
	"time"
)

type Token struct {
	AccessToken string `json:"accessToken"`
	ExpiresAt   int64  `json:"expiresAt"`
}

type TokenPayload struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	Password string `json:"password"`
	Verified bool   `json:"verified"`
	Key      string
}

type myClaims struct {
	Payload TokenPayload `json:"payload"`
	jwt.RegisteredClaims
}

func GenerateToken(data TokenPayload, tokenExpried int, secret string) (*Token, error) {
	expiresAt := jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(tokenExpried)))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, myClaims{
		data,
		jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			// Token được tạo khi nào
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ID:       fmt.Sprintf("%d", time.Now().UnixNano()),
		},
	})
	accessToken, _ := token.SignedString([]byte(secret))
	return &Token{
		accessToken,
		expiresAt.Unix(),
	}, nil
}
func ValidateJWT(accessToken string, cfg *config.Config) (*TokenPayload, error) {
	token, err := jwt.ParseWithClaims(accessToken, &myClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Services.Secret), nil
	})
	if err != nil {
		return nil, ErrTokenIsExpired
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*myClaims)
	if !ok {
		return nil, err
	}

	return &claims.Payload, nil
}

var (
	ErrTokenIsExpired = common.ErrUnauthorized(
		errors.New("token is expired"),
	)

	ErrTokenNotFound = common.ErrUnauthorized(
		errors.New("token not found"),
	)

	ErrEncodingToken = common.ErrUnauthorized(
		errors.New("error encoding token"),
	)

	ErrInvalidToken = common.ErrUnauthorized(
		errors.New("invalid token"),
	)
)
