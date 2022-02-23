package service

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"user/domain/model"
	"user/domain/repository"
)

////初始化数据表
//InitTable() error
////根据用户名称查找用户信息
//FindUserByName(string) (*model.User, error)
////根据用户ID查找用户信息
//FindUserByID(int64) (*model.User, error)
type IUserService interface {
	FindUserByName(string) (*model.User, error)
	FindUserByID(userID int64) (*model.User, error)
	AddUser(user *model.User) (int64, error)
	CheckPwd(userName string, pwd string) (isOk bool, err error)
}

//创建实例
func NewUserDataService(userRepository repository.IUserRepository) IUserService {
	return &UserDataService{UserRepository: userRepository}
}

type UserDataService struct {
	UserRepository repository.IUserRepository
}

func (u *UserDataService) AddUser(user *model.User) (int64, error) {
	return u.UserRepository.CreateUser(user)
}

func (u *UserDataService) FindUserByID(userID int64) (*model.User, error) {
	return u.UserRepository.FindUserByID(userID)
}

func (u *UserDataService) FindUserByName(name string) (*model.User, error) {
	return u.UserRepository.FindUserByName(name)
}

//检验密码【方法】
func (u *UserDataService) CheckPwd(userName string, pwd string) (isOk bool, err error) {
	//1.根据用户名查找用户信息
	user, err := u.UserRepository.FindUserByName(userName)
	if err != nil {
		return false, err
	}
	return valiatePassword(pwd, user.HashPassword)
}

//验证用户密码
func valiatePassword(userPassword string, hashed string) (isOk bool, err error) {
	if err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(userPassword)); err != nil {
		return false, errors.New("密码比对错误")
	}
	return true, nil
}

//加密用户密码
func GeneratePassword(userPassword string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
}
