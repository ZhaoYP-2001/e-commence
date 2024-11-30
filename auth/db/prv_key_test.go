package db

import (
	"e_commerce/auth/config"
	"errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestPrvKey(t *testing.T) {
	InitDB(config.MySQLDSN)
	user_id := "test_user"
	originPrvKey := "test_prv_key"

	err := CreatePrvKey(user_id, originPrvKey)
	if err != nil {
		t.Error(err)
	}
	outputPrvKey, err := QueryPrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, originPrvKey, outputPrvKey)

	createNewPrvKey := "test_create_prv_key"
	err = CreatePrvKey(user_id, createNewPrvKey)
	if err != nil {
		t.Error(err)
	}
	outputPrvKey, err = QueryPrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, createNewPrvKey, outputPrvKey)

	updateNewPrvKey := "test_update_prv_key"
	err = UpdatePrvKey(user_id, updateNewPrvKey)
	if err != nil {
		t.Error(err)
	}
	outputPrvKey, err = QueryPrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, updateNewPrvKey, outputPrvKey)

	err = DeletePrvKey(user_id)
	if err != nil {
		t.Error(err)
	}
	_, err = QueryPrvKey(user_id)
	assert.Equal(t, errors.Is(err, gorm.ErrRecordNotFound), true)
}
