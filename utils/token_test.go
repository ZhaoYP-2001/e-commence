package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRSA(t *testing.T) {
	prvKey, err := GenRsaKey(2048)
	if err != nil {
		t.Error(err)
	}

	pubKeyStr, err := GetRSAPubKeyStr(prvKey)
	if err != nil {
		t.Error(err)
	}

	pubKey, err := ParseRSAPubKey(pubKeyStr)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, &prvKey.PublicKey, pubKey)
}

func TestToken(t *testing.T) {
	user_id := "test_user"
	prvKey, err := GenRsaKey(2048)
	if err != nil {
		t.Error(err)
	}

	tokenStr, err := GenerateToken(prvKey, user_id)
	if err != nil {
		t.Error(err)
	}

	output_user_id, err := VerifyToken(&prvKey.PublicKey, tokenStr)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, user_id, output_user_id)
}
