package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lightsaid/short-net/util"
)

type Payload struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func NewPayload(userID uint, duration time.Duration) *Payload {
	var jwtID = fmt.Sprintf("%s%d%d", util.RandomString(10), time.Now().UnixMicro(), util.RandomInt(1, 10000))
	payload := &Payload{
		userID,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        jwtID,
			// Issuer:    "test",
			// Subject:   "somebody",
			// Audience:  []string{"somebody_else"},
		},
	}

	return payload
}
