package models

import (
	"fmt"
	"main/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	_DB *gorm.DB
)

func DB() *gorm.DB {
	return _DB
}

// ErrMissingRequiredField 缺少必要字段
type ErrMissingRequiredField struct {
	Field string // 缺少的字段
}

func (e ErrMissingRequiredField) Error() string {
	return fmt.Sprintf("missing required field %s", e.Field)
}

// ErrAlreadyExists 违反唯一性约束
type ErrAlreadyExists struct {
	Field string
	Value string
}

func (e ErrAlreadyExists) Error() string {
	return fmt.Sprintf("field %s with value %s already exists", e.Field, e.Value)
}

// ErrNotFound 未找到
type ErrNotFound struct {
	Model string
	Key   string
	Value string
}

func (e ErrNotFound) Error() string {
	return fmt.Sprintf("Cannot find %s with %s: %s", e.Model, e.Key, e.Value)
}

// Init 初始化数据库连接
//
//	@return error
func Init() error {
	db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect database: %v", err)
	}

	_DB = db

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&Favorite{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Follow{})

	return nil
}
