package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

// 定义Token错误类型
var (
	ErrInvalidToken = errors.New("令牌无效")
)

type Maker interface {
	GenToken(userID uint, duration time.Duration) (string, *Payload, error)
	VerifyToken(token string) (*Payload, error)
}

type TokenMaker struct {
	secretKey string
}

func NewTokenMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("tokenMaker 密钥长度必须大于或等于 %d", minSecretKeySize)
	}
	return &TokenMaker{secretKey}, nil
}

// GenToken 生成 token
func (maker *TokenMaker) GenToken(userID uint, duration time.Duration) (string, *Payload, error) {
	payload := NewPayload(userID, duration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	ss, err := token.SignedString([]byte(maker.secretKey))

	return ss, payload, err
}

// VerifyToken 验证 token
func (maker *TokenMaker) VerifyToken(tokenString string) (*Payload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Payload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(maker.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Payload); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}
