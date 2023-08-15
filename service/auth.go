package service

import (
	"main/config"
	"main/models"
	"main/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// ErrPasswordIncorrect 密码错误
type ErrPasswordIncorrect struct{}

func (e ErrPasswordIncorrect) Error() string {
	return "password incorrect"
}

// ErrInvalidToken 无效的token
type ErrInvalidToken struct{}

func (e ErrInvalidToken) Error() string {
	return "invalid token"
}

// ErrTokenExpired token过期
type ErrTokenExpired struct{}

func (e ErrTokenExpired) Error() string {
	return "token expired"
}

// Authenticate 验证用户
//
// authenticates a user with the given name and password.
// It returns the user ID if authentication is successful, or an error if authentication fails.
//
//	@param name
//	@param password
//	@return id
//	@return err
func Authenticate(name, password string) (id int64, err error) {
	user, err := models.UserDao().GetByName(name)
	if err != nil {
		return -1, err
	}
	temp := utils.Hash(password + user.Salt)
	if user.Password != temp {
		return -1, ErrPasswordIncorrect{}
	}
	return user.Id, nil
}

// GenerateToken 生成JWT Token
//
// generates a JWT token for the given user.
// It takes a pointer to a User struct as input and returns a string token and an error.
// The token contains the user's name, ID, and expiration time.
// The function uses the JWT signing method HS256 and the JWT secret from the config package.
//
//	@param user *User
//	@return string
//	@return error
func GenerateToken(user *models.User) (string, error) {
	currentTime := time.Now().Unix()
	expireTime := currentTime + config.ExpireTime
	claims := jwt.StandardClaims{
		Audience:  user.Name,
		ExpiresAt: expireTime,
		Id:        strconv.FormatInt(user.Id, 10),
		IssuedAt:  currentTime,
		Issuer:    "dy-svc",
		NotBefore: currentTime,
		Subject:   "login",
	}

	jwtSecret := []byte(config.JWTSecret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return token, nil
}

// verifyToken 验证JWT Token
//
// verifies the given JWT token and returns the token claims if valid.
// If the token is expired or invalid, an respective error is returned.
//
//	@param token
//	@return *jwt.StandardClaims
//	@return error
func verifyToken(token string) (*jwt.StandardClaims, error) {
	jwtSecret := []byte(config.JWTSecret)
	tokenClaims, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		// if err is start with "token is expired by", return ErrTokenExpired
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, ErrTokenExpired{}
		}
		return nil, ErrInvalidToken{}
	}
	claims, ok := tokenClaims.Claims.(*jwt.StandardClaims)
	if !ok || !tokenClaims.Valid {
		return nil, ErrInvalidToken{}
	}
	return claims, nil
}

// AuthenticateToken 验证JWT Token
//
// takes a JWT token as input and returns the user ID if the token is valid.
// respective errors are returned if the token is invalid or expired.
//
//	@param token
//	@return id
//	@return error
func AuthenticateToken(token string) (id int64, err error) {
	claims, err := verifyToken(token)
	if err != nil {
		return -1, err
	}
	uid, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return -1, ErrInvalidToken{}
	}
	return uid, nil
}
