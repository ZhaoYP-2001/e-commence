package service

import (
	"crypto/md5"
	"crypto/rsa"
	"e_commerce/auth/config"
	"e_commerce/auth/db"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

const (
	AccessToken  = "access_token"
	RefreshToken = "refresh_token"
)

// 参考：https://juejin.cn/post/7123943950616018958

func GenerateToken(prvKey *rsa.PrivateKey) (string, string, error) {
	claims := make(jwt.MapClaims)
	claims["type"] = AccessToken
	now := time.Now().UTC()
	claims["exp"] = now.Add(config.AccessTokenExpTime).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessTokenStr, err := accessToken.SignedString(prvKey)
	if err != nil {
		return "", "", fmt.Errorf("generate token error: %w", err)
	}

	claims["type"] = RefreshToken
	claims["exp"] = now.Add(config.RefreshTokenExpTime).Unix()
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	refreshTokenStr, err := refreshToken.SignedString(prvKey)
	if err != nil {
		return "", "", fmt.Errorf("generate token error: %w", err)
	}

	err = db.SetToken(getMD5Hash(refreshTokenStr))
	if err != nil {
		return "", "", fmt.Errorf("generate token error: %w", err)
	}

	return accessTokenStr, refreshTokenStr, nil
}

func VerifyAccessToken(pubKey *rsa.PublicKey, tokenStr string) error {
	claims, err := parseToken(pubKey, tokenStr)
	if err != nil {
		return fmt.Errorf("verify access token error: %w", err)
	}

	if tokenType, ok := claims["type"]; ok && tokenType.(string) == AccessToken {
		return nil
	} else {
		return fmt.Errorf("verify access token error: unsupproted token: %v", tokenStr)
	}
}

func VerifyRefreshToken(pubKey *rsa.PublicKey, tokenStr string) error {
	claims, err := parseToken(pubKey, tokenStr)
	if err != nil {
		return fmt.Errorf("verify refresh token error: %w", err)
	}

	if tokenType, ok := claims["type"]; !ok || tokenType.(string) != RefreshToken {
		return fmt.Errorf("verify refresh token error: unsupproted token: %v", tokenStr)
	}

	hashValue := getMD5Hash(tokenStr)
	_, err = db.GetToken(hashValue)
	if err != nil {
		return fmt.Errorf("verify refresh token error: %v", tokenStr)
	}

	err = db.DelToken(hashValue)
	if err != nil {
		return fmt.Errorf("verify refresh token error: %v", tokenStr)
	}

	return nil
}

func getMD5Hash(tokenStr string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(tokenStr)))
}

func parseToken(pubKey *rsa.PublicKey, tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %w", jwtToken.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token error: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, fmt.Errorf("invalid token: %v", tokenStr)
	}
}
