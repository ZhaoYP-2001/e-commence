package service

import (
	"bou.ke/monkey"
	"e_commerce/auth/db"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestPrvKey(t *testing.T) {
	storedPrvKeyMap := sync.Map{}
	monkey.Patch(db.CreatePrvKey, func(userId string, prvKey string) error {
		storedPrvKeyMap.Store(userId, prvKey)
		return nil
	})
	defer monkey.Unpatch(db.CreatePrvKey)

	monkey.Patch(db.QueryPrvKey, func(userId string) (string, error) {
		if prvKey, ok := storedPrvKeyMap.Load(userId); ok {
			return prvKey.(string), nil
		} else {
			return "", fmt.Errorf("no private key for user %w", userId)
		}
	})
	defer monkey.Unpatch(db.QueryPrvKey)

	user_id := "test_user"
	expectedPrvKey, err := NewPrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	outputPrvKey, err := LoadPrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, expectedPrvKey, outputPrvKey)
}
