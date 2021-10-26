package services

import (
	"fmt"
	"github.com/arrayJY/lite-im-server/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strings"
	"time"
)

var Key = []byte("AllYourBase")

func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Need authorization."})
			return
		}
		s := strings.Split(authorization, " ")
		if s[0] != "Bearer" || len(s) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bad authorization."})
			return
		}
		tokenString := s[1]
		token, err := jwt.ParseWithClaims(tokenString, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v.", token.Header["alg"])
			}
			return Key, nil
		})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		if claims, ok := token.Claims.(*models.AccessTokenClaims); ok && token.Valid {
			c.Set("AuthorizationID", claims.Id)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bad token."})
			return
		}
	}
}

const expiresSeconds int64 = 1200 //20min

func signToken() (*models.TokenResponse, error) {
	now := time.Now().Unix()
	tokenExpiresAt := now + expiresSeconds
	refreshTokenExpiresAt := now + expiresSeconds*10
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.AccessTokenClaims{
		Id: 0,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiresAt,
			Issuer:    "test",
		},
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiresAt,
			Issuer:    "test",
		},
	})
	accessTokenString, err1 := accessToken.SignedString(Key)
	if err1 != nil {
		return nil, err1
	}
	refreshTokenString, err2 := refreshToken.SignedString(Key)
	if err1 != nil {
		return nil, err2
	}
	return &models.TokenResponse{
		AccessToken:  accessTokenString,
		Expires:      tokenExpiresAt,
		RefreshToken: refreshTokenString,
	}, nil
}

func CreateToken(c *gin.Context) {
	var loginInfo models.LoginRequest
	if err := c.BindJSON(&loginInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO: authorize name and password
	if token, err := signToken(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusOK, token)
	}
}

func RefreshToken(c *gin.Context) {
	var body models.RefreshTokenRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	token, err := jwt.ParseWithClaims(body.RefreshToken, &models.AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("bad token: unexpected signing method %v", token.Header["alg"])
		}
		return Key, nil
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
	if token, err := signToken(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, token)
	}
}
