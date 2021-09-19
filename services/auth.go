package services

import (
	"github.com/arrayJY/go-im-server/models"
	"github.com/emicklei/go-restful/v3"
	"github.com/golang-jwt/jwt"
	"time"
)

func AuthService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/auth").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	// Sign in.
	ws.Route(ws.
		POST("/token").
		Filter(auth).
		To(createToken))

	// Refresh token.
	ws.Route(ws.
		PUT("/token").
		Filter(authToken).
		To(createToken))
	return ws
}

const expiresSeconds int64 = 1200 //20min
var mySigningKey = []byte("AllYourBase")

func createToken(_ *restful.Request, resp *restful.Response) {
	now := time.Now().Unix()
	tokenExpiresAt := now + expiresSeconds
	refreshTokenExpiresAt := now + expiresSeconds*10
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.AccessTokenClaims{
		Id: "", //TODO: put id into claims
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiresAt,
			Issuer:    "test",
		},
	})
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, models.RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: refreshTokenExpiresAt, // 200min expires
			Issuer:    "test",
		},
	})

	tokenString, tokenErr := token.SignedString(mySigningKey)
	refreshTokenString, refreshTokenErr := refreshToken.SignedString(mySigningKey)
	if tokenErr != nil || refreshTokenErr != nil {
		return
	}
	resp.WriteAsJson(models.TokenResponse{
		Token:        tokenString,
		Expires:      tokenExpiresAt,
		RefreshToken: refreshTokenString,
	})
}

func auth(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//TODO: authorization
	chain.ProcessFilter(req, resp)
}

func authToken(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	//TODO: authorization
	chain.ProcessFilter(req, resp)
}
