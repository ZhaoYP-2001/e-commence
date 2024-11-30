package db

import (
	"e_commerce/auth/config"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToken(t *testing.T) {
	InitRedis(config.RedisAddr, config.RedisPassword)
	tokenHash := "hash value"
	err := SetToken(tokenHash)
	if err != nil {
		t.Error(err)
	}

	_, err = GetToken(tokenHash)
	assert.Equal(t, nil, err)

	err = DelToken(tokenHash)
	assert.Equal(t, nil, err)

	_, err = GetToken(tokenHash)
	assert.Equal(t, errors.Is(err, redis.Nil), true)
}
