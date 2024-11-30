package config

import "time"

const (
	AccessTokenExpTime  = time.Hour * 24
	RefreshTokenExpTime = time.Hour * 24 * 3
)
const (
	RSAKeyBits = 2048
)

const (
	ServerAddr = "localhost:9091"
)

const (
	MySQLDSN      = "root:123456@tcp(localhost:3306)/auth?charset=utf8mb4&parseTime=True&loc=Local"
	RedisAddr     = "localhost:6379"
	RedisPassword = ""
)
