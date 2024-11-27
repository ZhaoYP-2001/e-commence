package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	"time"
)

// 参考：https://juejin.cn/post/7123943950616018958

func GenerateToken(prvKey *rsa.PrivateKey, user_id string) (string, error) {
	claims := make(jwt.MapClaims)
	claims["user_id"] = user_id
	now := time.Now().UTC()
	claims["exp"] = now.Add(time.Hour * 24).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(prvKey)
	if err != nil {
		return "", fmt.Errorf("generate token error: %s", err.Error())
	}
	return token, nil
}

func VerifyToken(pubKey *rsa.PublicKey, tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}
		return pubKey, nil
	})

	if err != nil {
		return "", fmt.Errorf("parse token error: %s", err.Error())
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["user_id"].(string), nil
	} else {
		return "", fmt.Errorf("invalid token: %v", tokenStr)
	}
}

// 参考文档：https://cloud.tencent.com/developer/article/2339220

func GenRsaKey(bits int) (*rsa.PrivateKey, error) {
	// Generates private key.
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("generate rsa key error: %s", err.Error())
	}
	return privateKey, nil
}

func GetRSAPubKeyStr(privateKey *rsa.PrivateKey) (string, error) {

	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", fmt.Errorf("convert rsa key error: %s", err.Error())
	}
	block := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: derPkix,
	}
	pubkeyStr := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(block))
	return pubkeyStr, nil
}

func ParseRSAPubKey(pubkeyStr string) (*rsa.PublicKey, error) {
	pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyStr)
	if err != nil {
		return nil, fmt.Errorf("parse public key error: %s", err.Error())
	}
	block, _ := pem.Decode(pubkeyBytes)
	if block == nil {
		return nil, fmt.Errorf("parse public key error: %s", err.Error())
	}
	pubkey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse public key error: %s", err.Error())
	}
	return pubkey.(*rsa.PublicKey), nil
}
