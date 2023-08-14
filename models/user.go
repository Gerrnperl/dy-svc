package models

import (
	"main/config"
	"main/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Id              int64  `json:"id,omitempty" gorm:"primarykey"`
	Name            string `json:"name,omitempty"`
	Password        string `json:"password,omitempty"`
	Salt            string `json:"salt,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	Signature       string `json:"signature,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

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

// AddUser 添加用户
//
//	@param user *User 用户, 必填字段: Name, Password
//	@return *User
//	@return error
func AddUser(user *User) (*User, error) {
	// 判断必填字段是否为空
	if user.Name == "" {
		return nil, ErrMissingRequiredField{"name"}
	}
	if user.Password == "" {
		return nil, ErrMissingRequiredField{"password"}
	}
	// 判断用户名是否已存在
	var count int64
	DB().Model(&User{}).Where("name = ?", user.Name).Count(&count)
	if count > 0 {
		return nil, ErrAlreadyExists{"name", user.Name}
	}

	pwd, salt := utils.HashWithSalt(user.Password)

	newUser := User{
		Name:     user.Name,
		Password: pwd,
		Salt:     salt,
	}

	result := DB().Create(&newUser)
	if result.Error != nil {
		return nil, result.Error
	}
	return &newUser, nil
}

// GetUserByName 根据用户名获取用户
//
//	@param name
//	@return *User
//	@return error
func GetUserByName(name string) (*User, error) {
	var user User
	result := DB().Where("name = ?", name).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{
				"user",
				"name",
				name,
			}
		}
	}
	return &user, nil
}

// GetUserById 根据用户ID获取用户
//
//	@param id
//	@return *User
//	@return error
func GetUserById(id int64) (*User, error) {
	var user User
	result := DB().Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, ErrNotFound{
				"user",
				"id",
				strconv.FormatInt(id, 10),
			}
		}
		return nil, result.Error
	}
	return &user, nil
}

// Authenticate 验证用户
//
//	@param name string 用户名
//	@param password string 密码
//	@return *User
//	@return error
func Authenticate(name, password string) (*User, error) {
	user, err := GetUserByName(name)
	if err != nil {
		return nil, err
	}
	temp := utils.Hash(password + user.Salt)
	if user.Password != temp {
		return nil, ErrPasswordIncorrect{}
	}
	return user, nil
}

// GenerateToken 生成JWT Token
//
//	@param user *User
//	@return string
//	@return error
func GenerateToken(user *User) (string, error) {
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
//	@param token
//	@return *User
//	@return error
func AuthenticateToken(token string) (*User, error) {
	claims, err := verifyToken(token)
	if err != nil {
		return nil, err
	}
	id, err := strconv.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return nil, ErrInvalidToken{}
	}
	user, err := GetUserById(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
