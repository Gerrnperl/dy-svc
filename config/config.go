package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	DSN        string
	Address    string
	Port       string
	ExpireTime int64 = 60 * 60 * 24 // 1 day
	JWTSecret  string
)

// readEnv 从环境变量中读取配置, 如果不存在则报错并退出
//
//	@int key
//	@return string
//	@return error
func readEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("environment variable %s is required", key)
	}
	return value, nil
}

// readEnvWithDefault 从环境变量中读取配置, 如果不存在则使用默认值
//
//	@param key
//	@param defaultValue
//	@return string
func readEnvWithDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

// Init 从环境变量中读取配置
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	DBUser, _ := readEnv("DB_USER")
	DBPwd, _ := readEnv("DB_PWD")
	DBHost, _ := readEnv("DB_HOST")
	DBPort, _ := readEnv("DB_PORT")
	DBName, _ := readEnv("DB_NAME")
	DBConfig := readEnvWithDefault("DB_CONFIG", "charset=utf8mb4&parseTime=True&loc=Local")
	DSN = DBUser + ":" + DBPwd + "@tcp(" + DBHost + ":" + DBPort + ")/" + DBName + "?" + DBConfig

	Address = readEnvWithDefault("ADDR", "")
	Port = readEnvWithDefault("PORT", "8080")

	expireStr := readEnvWithDefault("EXPIRE_TIME", "86400")
	expire, err := strconv.Atoi(expireStr)
	if err != nil {
		log.Fatalf("invalid expire: %s", expireStr)
	}
	ExpireTime = int64(expire)

	JWTSecret, _ = readEnv("JWT_SECRET")
}
