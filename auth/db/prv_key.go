package db

import (
	"errors"
	"fmt"
	"gorm.io/gorm/clause"
	"gorm.io/plugin/soft_delete"
	"time"
)

type PreKey struct {
	UserId    string `gorm:"primaryKey"`
	PrvKey    string
	CreatedAt time.Time
	DeletedAt time.Time             `gorm:"default:Null"`
	IsDel     soft_delete.DeletedAt `gorm:"softDelete:flag,DeletedAtField:DeletedAt"`
}

func (prvKey *PreKey) TableName() string {
	return "private_key"
}

func CreatePrvKey(userId string, prvKey string) error {
	if userId == "" || prvKey == "" {
		return errors.New("save private key to database error: user-id and private key should not be empty")
	}
	result := db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&PreKey{UserId: userId, PrvKey: prvKey})
	if result.Error != nil {
		return fmt.Errorf("save private key to database error: %w", result.Error)
	}
	return nil
}

func DeletePrvKey(userId string) error {
	result := db.Delete(&PreKey{UserId: userId})
	if result.Error != nil {
		return fmt.Errorf("delete private key from database error: %w", result.Error)
	}
	return nil
}

func UpdatePrvKey(userId string, prvKey string) error {
	if userId == "" || prvKey == "" {
		return errors.New("save private key to database error: user-id and private key should not be empty")
	}
	result := db.Model(&PreKey{UserId: userId}).Updates(&PreKey{PrvKey: prvKey})
	if result.Error != nil {
		return fmt.Errorf("save private key to database error: %w", result.Error)
	}
	return nil
}

func QueryPrvKey(userId string) (string, error) {
	prvKey := &PreKey{UserId: userId}
	result := db.First(prvKey)
	if result.Error != nil {
		return "", fmt.Errorf("load private key from database error: %w", result.Error)
	}
	return prvKey.PrvKey, nil
}
