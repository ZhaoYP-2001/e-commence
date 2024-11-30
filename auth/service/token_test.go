package service

import (
	"bou.ke/monkey"
	"e_commerce/auth/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestAccessToken(t *testing.T) {
	monkey.Patch(db.SetToken, func(_ string) error {
		return nil
	})
	defer monkey.Unpatch(db.SetToken)

	prvKey, err := genRsaKey(2048)
	if err != nil {
		t.Error(err)
	}

	accessToken, _, err := GenerateToken(prvKey)
	if err != nil {
		t.Error(err)
	}

	err = VerifyAccessToken(&prvKey.PublicKey, accessToken)
	assert.Equal(t, nil, err)
}

func TestRefreshToken(t *testing.T) {
	storedTokenMap := sync.Map{}
	monkey.Patch(db.SetToken, func(tokenHash string) error {
		storedTokenMap.Store(tokenHash, 1)
		return nil
	})
	defer monkey.Unpatch(db.SetToken)

	monkey.Patch(db.GetToken, func(tokenHash string) (int, error) {
		if token, ok := storedTokenMap.Load(tokenHash); ok {
			return token.(int), nil
		} else {
			return 0, fmt.Errorf("no token for hash value %w", tokenHash)
		}
	})
	defer monkey.Unpatch(db.GetToken)

	monkey.Patch(db.DelToken, func(tokenHash string) error {
		storedTokenMap.Delete(tokenHash)
		return nil
	})
	defer monkey.Unpatch(db.DelToken)

	prvKey, err := genRsaKey(2048)
	if err != nil {
		t.Error(err)
	}

	_, refreshToken, err := GenerateToken(prvKey)
	if err != nil {
		t.Error(err)
	}

	err = VerifyRefreshToken(&prvKey.PublicKey, refreshToken)
	assert.Equal(t, nil, err)

	tokenCount := 0
	storedTokenMap.Range(func(key, value interface{}) bool {
		tokenCount++
		return true // 继续迭代
	})
	assert.Equal(t, 0, tokenCount)
}
