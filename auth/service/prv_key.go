package service

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"e_commerce/auth/config"
	"e_commerce/auth/db"
	"encoding/pem"
	"errors"
	"fmt"
)

// 参考：https://cloud.tencent.com/developer/article/2339220

func genRsaKey(bits int) (*rsa.PrivateKey, error) {
	prvkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("generate rsa key error: %w", err)
	}
	return prvkey, nil
}

func encodePrvKey(prvkey *rsa.PrivateKey) string {
	derStream := x509.MarshalPKCS1PrivateKey(prvkey)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	return string(pem.EncodeToMemory(block))
}

func parsePrvKey(prvKeyStr string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(prvKeyStr))
	if block == nil {
		return nil, errors.New("parse private key error")
	}
	prvKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key error: %w", err)
	}
	return prvKey, nil
}

func SavePrvKey(userId string, prvKey *rsa.PrivateKey) error {
	if userId == "" {
		return errors.New("save private key error: user-id and private key should not be empty")
	}
	if err := prvKey.Validate(); err != nil {
		return fmt.Errorf("save private key error: %w", err)
	}
	prvKeyStr := encodePrvKey(prvKey)
	err := db.UpdatePrvKey(userId, prvKeyStr)
	if err != nil {
		return fmt.Errorf("save private key error: %w", err)
	}
	return nil
}

func LoadPrvKey(userId string) (*rsa.PrivateKey, error) {
	prvKeyStr, err := db.QueryPrvKey(userId)
	if err != nil {
		return nil, fmt.Errorf("load private key error: %w", err)
	}
	prvKey, err := parsePrvKey(prvKeyStr)
	if err != nil {
		return nil, fmt.Errorf("load private key error: %w", err)
	}
	return prvKey, nil
}

func NewPrvKey(userId string) (*rsa.PrivateKey, error) {
	prvKey, err := genRsaKey(config.RSAKeyBits)
	if err != nil {
		return nil, fmt.Errorf("generate private key error: %w", err)
	}
	err = db.CreatePrvKey(userId, encodePrvKey(prvKey))
	if err != nil {
		return nil, fmt.Errorf("generate private key error: %w", err)
	}
	return prvKey, nil
}
