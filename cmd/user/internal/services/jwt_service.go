package services

import (
	"github.com/tricong1998/go-ecom/cmd/user/internal/config"
	"github.com/tricong1998/go-ecom/pkg/token"
)

type JwtService struct {
	tokenMaker token.Maker
	authConfig config.AuthConfig
}

type IJwtService interface {
	CreateToken(username string, userId uint, role string) (string, string, error)
	VerifyToken(token string) (string, error)
}

func NewJwtService(tokenMaker token.Maker, authConfig config.AuthConfig) *JwtService {
	return &JwtService{tokenMaker, authConfig}
}

func (jwtService *JwtService) CreateToken(username string, userId uint, role string) (string, string, error) {
	accessToken, _, err := jwtService.tokenMaker.CreateToken(username, userId, jwtService.authConfig.AccessTokenDuration, role)
	if err != nil {
		return "", "", err
	}
	refreshToken, _, err := jwtService.tokenMaker.CreateToken(username, userId, jwtService.authConfig.RefreshTokenDuration, role)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (jwtService *JwtService) VerifyToken(token string) (string, error) {
	payload, err := jwtService.tokenMaker.VerifyToken(token)
	if err != nil {
		return "", err
	}
	return payload.Username, nil
}
