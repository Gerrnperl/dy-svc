package service

import (
	"main/config"
	"main/models"
	"reflect"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey"
	"github.com/dgrijalva/jwt-go"
)

var mockUser = &models.User{
	Id:              1,
	Name:            "test",
	Password:        "14f721c943c2bc7aae6b031126fbe3f9",
	Salt:            "hello123salt",
	FollowCount:     1,
	FollowerCount:   1,
	Avatar:          "testavatar",
	BackgroundImage: "testbgimg",
	TotalFavorited:  1,
	WorkCount:       1,
	FavoriteCount:   1,
	Signature:       "testsignature",
}

func TestVerifyToken(t *testing.T) {
	token, err := GenerateToken(&models.User{
		Id:   1,
		Name: "test",
	})
	if err != nil {
		t.Error(err)
	}
	claims, err := verifyToken(token)
	if err != nil {
		t.Error(err)
	}
	if claims.Audience != "test" {
		t.Error("audience error")
	}
	if claims.Id != "1" {
		t.Error("id error")
	}
}

func TestVerifyTokenExpired(t *testing.T) {
	// 生成过期的 JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.StandardClaims{
		Id:        "1",
		Audience:  "test",
		ExpiresAt: time.Now().Unix() - 1000,
	})
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
	if err != nil {
		t.Error(err)
	}

	// 调用 verifyToken 函数并检查其返回值是否为预期的错误
	_, err = verifyToken(tokenString)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrTokenExpired); !ok {
		t.Error("expected ErrTokenExpired, but got", err)
	}
}

func TestAuthenticateToken(t *testing.T) {

	// 生成虚假的 JWT Token
	token, err := GenerateToken(&models.User{
		Id:   1,
		Name: "test",
	})
	if err != nil {
		t.Error(err)
	}

	// 调用 AuthenticateToken 函数并检查其返回值是否为预期的虚假用户对象
	id, err := AuthenticateToken(token)
	if err != nil {
		t.Error(err)
	}

	if id != 1 {
		t.Error("id not equal")
	}
}

func TestAuthenticateTokenInvalidToken(t *testing.T) {
	// 调用 AuthenticateToken 函数并检查其返回值是否为预期的错误
	_, err := AuthenticateToken("invalidtoken")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrInvalidToken); !ok {
		t.Error("expected ErrInvalidToken, but got", err)
	}
}

func TestAuthenticate(t *testing.T) {
	// 模拟数据库查询结果
	gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "GetByName", func(*models.UserDaoStruct, string) (*models.User, error) {
		return mockUser, nil
	})

	// 调用 Authenticate 函数并检查其返回值是否为预期的虚假用户对象
	id, err := Authenticate("test", "testpwd")
	if err != nil {
		t.Error(err)
	}

	if id != mockUser.Id {
		t.Error("user id not equal")
	}
}

func TestAuthenticateUserNotFound(t *testing.T) {
	// 替换models.UserDao().GetByName(name)
	gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "GetByName", func(*models.UserDaoStruct, string) (*models.User, error) {
		return nil, models.ErrNotFound{Model: "user", Key: "name", Value: "test"}
	})

	// 调用 Authenticate 函数并检查其返回值是否为预期的错误
	_, err := Authenticate("test", "testpwd")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(models.ErrNotFound); !ok {
		t.Error("expected ErrNotFound, but got", err)
	}
}

func TestAuthenticateIncorrectPassword(t *testing.T) {
	// 模拟数据库查询结果
	gomonkey.ApplyMethod(reflect.TypeOf(models.UserDao()), "GetByName", func(*models.UserDaoStruct, string) (*models.User, error) {
		return mockUser, nil
	})

	// 调用 Authenticate 函数并检查其返回值是否为预期的错误
	_, err := Authenticate("test", "wrongpwd")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrPasswordIncorrect); !ok {
		t.Error("expected ErrPasswordIncorrect, but got", err)
	}
}
