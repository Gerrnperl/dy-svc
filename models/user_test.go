package models

import (
	"main/config"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/dgrijalva/jwt-go"
	"gorm.io/gorm"
)

func TestVerifyToken(t *testing.T) {
	token, err := GenerateToken(&User{
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
	token, err := GenerateToken(&User{
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

var mockUser = &User{
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

func TestGetUserById(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE id = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "salt", "follow_count", "follower_count", "avatar", "background_image", "total_favorited", "work_count", "favorite_count", "signature"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.Password, mockUser.Salt, mockUser.FollowCount, mockUser.FollowerCount, mockUser.Avatar, mockUser.BackgroundImage, mockUser.TotalFavorited, mockUser.WorkCount, mockUser.FavoriteCount, mockUser.Signature),
	)

	// 调用 GetUserById 函数并检查其返回值是否为预期的虚假用户对象
	user, err := GetUserById(1)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(mockUser, user) {
		t.Error("user not equal")
	}
}

func TestGetUserByIdNotFound(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE id = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs(1).WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "salt", "follow_count", "follower_count", "avatar", "background_image", "total_favorited", "work_count", "favorite_count", "signature"}),
	)

	// 调用 GetUserById 函数并检查其返回值是否为预期的错误
	_, err := GetUserById(1)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrNotFound); !ok {
		t.Error("expected ErrNotFound, but got", err)
	}
}

func TestAddUser(t *testing.T) {
	user := &User{
		Name:     "test",
		Password: "testpwd",
	}

	mock.ExpectQuery("SELECT count(*) FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL").WithArgs(user.Name).WillReturnRows(
		sqlmock.NewRows([]string{"count(*)"}).AddRow(0),
	)

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO `user` (`created_at`,`updated_at`,`deleted_at`,`name`,`password`,`salt`,`follow_count`,`follower_count`,`avatar`,`background_image`,`total_favorited`,`work_count`,`favorite_count`,`signature`) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, user.Name, sqlmock.AnyArg(), sqlmock.AnyArg(), 0, 0, "", "", 0, 0, 0, "").WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// 调用 AddUser 函数并检查其返回值是否为预期的虚假用户对象
	newUser, err := AddUser(user)
	if err != nil {
		t.Error(err)
	}

	if newUser.Name != user.Name {
		t.Error("name not equal")
	}
	if newUser.Password == user.Password {
		t.Error("password not hashed")
	}
	if newUser.Salt == "" {
		t.Error("salt not generated")
	}
}

func TestAddUserEmptyName(t *testing.T) {
	user := &User{
		Name:     "",
		Password: "testpwd",
	}

	// 调用 AddUser 函数并检查其返回值是否为预期的错误
	_, err := AddUser(user)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrMissingRequiredField); !ok {
		t.Error("expected ErrMissingRequiredField, but got", err)
	}
}

func TestAddUserEmptyPassword(t *testing.T) {
	user := &User{
		Name:     "test",
		Password: "",
	}

	// 调用 AddUser 函数并检查其返回值是否为预期的错误
	_, err := AddUser(user)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrMissingRequiredField); !ok {
		t.Error("expected ErrMissingRequiredField, but got", err)
	}
}

func TestAddUserAlreadyExists(t *testing.T) {
	user := &User{
		Name:     "test",
		Password: "testpwd",
	}

	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT count(*) FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL").WithArgs(user.Name).WillReturnRows(
		sqlmock.NewRows([]string{"count(*)"}).AddRow(1),
	)

	// 调用 AddUser 函数并检查其返回值是否为预期的错误
	_, err := AddUser(user)
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrAlreadyExists); !ok {
		t.Error("expected ErrAlreadyExists, but got", err)
	}
}

func TestGetUserByName(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs("test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "salt", "follow_count", "follower_count", "avatar", "background_image", "total_favorited", "work_count", "favorite_count", "signature"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.Password, mockUser.Salt, mockUser.FollowCount, mockUser.FollowerCount, mockUser.Avatar, mockUser.BackgroundImage, mockUser.TotalFavorited, mockUser.WorkCount, mockUser.FavoriteCount, mockUser.Signature),
	)

	// 调用 GetUserByName 函数并检查其返回值是否为预期的虚假用户对象
	user, err := GetUserByName("test")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(mockUser, user) {
		t.Error("user not equal")
	}
}

func TestGetUserByNameNotFound(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs("test").WillReturnError(gorm.ErrRecordNotFound)

	// 调用 GetUserByName 函数并检查其返回值是否为预期的错误
	user, err := GetUserByName("test")
	if user != nil {
		t.Error("expected nil user, but got", user)
	}
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if err != err.(ErrNotFound) {
		t.Error("expected ErrNotFound, but got", err)
	}
}

func TestAuthenticate(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs("test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "salt", "follow_count", "follower_count", "avatar", "background_image", "total_favorited", "work_count", "favorite_count", "signature"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.Password, mockUser.Salt, mockUser.FollowCount, mockUser.FollowerCount, mockUser.Avatar, mockUser.BackgroundImage, mockUser.TotalFavorited, mockUser.WorkCount, mockUser.FavoriteCount, mockUser.Signature),
	)

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
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs("test").WillReturnError(gorm.ErrRecordNotFound)

	// 调用 Authenticate 函数并检查其返回值是否为预期的错误
	_, err := Authenticate("test", "testpwd")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrNotFound); !ok {
		t.Error("expected ErrNotFound, but got", err)
	}
}

func TestAuthenticateIncorrectPassword(t *testing.T) {
	// 模拟数据库查询结果
	mock.ExpectQuery("SELECT * FROM `user` WHERE name = ? AND `user`.`deleted_at` IS NULL ORDER BY `user`.`id` LIMIT 1").WithArgs("test").WillReturnRows(
		sqlmock.NewRows([]string{"id", "name", "password", "salt", "follow_count", "follower_count", "avatar", "background_image", "total_favorited", "work_count", "favorite_count", "signature"}).
			AddRow(mockUser.Id, mockUser.Name, mockUser.Password, mockUser.Salt, mockUser.FollowCount, mockUser.FollowerCount, mockUser.Avatar, mockUser.BackgroundImage, mockUser.TotalFavorited, mockUser.WorkCount, mockUser.FavoriteCount, mockUser.Signature),
	)

	// 调用 Authenticate 函数并检查其返回值是否为预期的错误
	_, err := Authenticate("test", "wrongpwd")
	if err == nil {
		t.Error("expected error, but got nil")
	}
	if _, ok := err.(ErrPasswordIncorrect); !ok {
		t.Error("expected ErrPasswordIncorrect, but got", err)
	}
}
