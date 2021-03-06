package models

import "github.com/golang-jwt/jwt"

type AccessTokenClaims struct {
	Id uint `json:"id"`
	jwt.StandardClaims
}
type RefreshTokenClaims struct {
	jwt.StandardClaims
}

type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type ErrorResponse struct {
	Detail struct {
		Msg string `json:"msg"`
	} `json:"detail"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	Expires      int64  `json:"expires"`
	RefreshToken string `json:"refresh_token"`
}
