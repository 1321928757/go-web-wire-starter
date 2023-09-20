package domain

import "github.com/golang-jwt/jwt/v5"

const (
	AppGuardName = "internal"
)

type JwtUser interface {
	GetUid() string
}

// CustomClaims 自定义 Claims
type CustomClaims struct {
	Key string `json:"key,omitempty"`
	jwt.RegisteredClaims
}

type TokenOutPut struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
