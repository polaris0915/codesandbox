package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/polaris/codesandbox/settings"
)

type UserClaims struct {
	jwt.RegisteredClaims
	Identity    string `json:"identity"`
	UserAccount string `json:"userAccount"`
	UserRole    string `json:"userRole"`
}

func ValidateToken(tokenString string) (*UserClaims, error) {
	userClaims := &UserClaims{}
	var jwtKey = []byte(settings.RemoteConfig.JwtConfig.Key)
	token, err := jwt.ParseWithClaims(tokenString, userClaims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return userClaims, nil
}
