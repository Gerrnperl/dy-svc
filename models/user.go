package models

import (
	"main/utils"
	"strconv"
	"sync"

	"gorm.io/gorm"
)

type UserProfile struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	Signature       string `json:"signature,omitempty"`
}

type User struct {
	gorm.Model

	Id              int64  `json:"id,omitempty" gorm:"primarykey"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	BackgroundImage string `json:"background_image,omitempty"`
	TotalFavorited  int64  `json:"total_favorited,omitempty"`
	WorkCount       int64  `json:"work_count,omitempty"`
	FavoriteCount   int64  `json:"favorite_count,omitempty"`
	Signature       string `json:"signature,omitempty"`

	Password string `json:"password,omitempty"`
	Salt     string `json:"salt,omitempty"`
}

func (u *User) TableName() string {
	return "user"
}

var (
	_userDaoInstance *UserDaoStruct
	_userDaoOnce     sync.Once
)

type UserDaoStruct struct{}

// UserDao returns a singleton instance of userDaoStruct.
//
// It uses sync.Once to ensure that only one instance of userDaoStruct is created.
//
// The returned instance can be used to perform CRUD operations on the user data.
func UserDao() *UserDaoStruct {
	_userDaoOnce.Do(func() {
		_userDaoInstance = &UserDaoStruct{}
	})
	return _userDaoInstance
}

// Add 添加用户
//
// Add adds a new user to the database. It takes a pointer to a User struct as input and returns a pointer to the newly created User struct and an error (if any).
//
// If the required fields (name and password) are missing, it returns an ErrMissingRequiredField error.
//
// If the user with the same name already exists in the database, it returns an ErrAlreadyExists error.
//
// It hashes the password and generates a salt before storing it in the database.
//
//	@param user *User
//	@return *User
//	@return error
func (dao *UserDaoStruct) Add(user *User) (*User, error) {
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

// GetByName 根据用户名获取用户
//
// GetByName retrieves a user from the database by name. It takes a string representing the user's name as input and returns a pointer to the User struct and an error (if any).
//
// If the user with the specified name is not found in the database, it returns an ErrNotFound error.
//
//	@param name
//	@return *User
//	@return error
func (dao *UserDaoStruct) GetByName(name string) (*User, error) {
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

// GetById 根据用户ID获取用户
//
//	@param id
//	@return *User
//	@return error
func (dao *UserDaoStruct) GetById(id int64) (*User, error) {
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
