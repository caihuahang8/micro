package repository

import (
	"github.com/jinzhu/gorm"
	"user/domain/model"
)

type IUserRepository interface {
	//初始化数据表
	InitTable() error
	//根据用户名称查找用户信息
	FindUserByName(string) (*model.User, error)
	//根据用户ID查找用户信息
	FindUserByID(int64) (*model.User, error)
	CreateUser(*model.User) (int64, error)
}

//创建UserRepository
func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{mysqlDb: db}
}

type UserRepository struct {
	mysqlDb *gorm.DB
}

func (u *UserRepository) CreateUser(user *model.User) (int64, error) {
	return user.ID, u.mysqlDb.Create(user).Error
}

func (u *UserRepository) InitTable() error {
	return u.mysqlDb.CreateTable(&model.User{}).Error
}
func (UserRepository *UserRepository) FindUserByName(name string) (*model.User, error) {
	u2 := &model.User{}
	return u2, UserRepository.mysqlDb.Where("user_name = ?", name).Find(u2).Error
}

func (UserRepository *UserRepository) FindUserByID(userId int64) (*model.User, error) {
	user := &model.User{}
	return user, UserRepository.mysqlDb.Where("user_id = ?", userId).Find(user).Error
}
