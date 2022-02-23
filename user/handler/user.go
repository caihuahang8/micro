package handler

import (
	"context"
	"user/domain/model"
	"user/domain/service"
	user "user/proto"
)

type User struct {
	UserService service.IUserService
}

func (u *User) Login(ctx context.Context, request *user.UserLoginRequest, response *user.UserLoginResponse) error {
	_, err := u.UserService.CheckPwd(request.GetUserName(), request.GetPwd())
	if err != nil {
		return err
	}
	response.IsSuccess = true
	return nil
}

func (u *User) GetUserInfo(ctx context.Context, request *user.UserInfoRequest, response *user.UserInfoResponse) error {
	userInfo, err := u.UserService.FindUserByName(request.GetUserName())
	if err != nil {
		return err
	}
	response = UserForResponse(userInfo)
	return nil
}

/**
注册
*/
func (u *User) Register(ctx context.Context, request *user.UserRegisterRequest, response *user.UserRegisterResponse) error {
	user := &model.User{
		UserName:     request.UserName,
		FirstName:    request.FirstName,
		HashPassword: request.Pwd,
	}
	_, err := u.UserService.AddUser(user)
	if err != nil {
		return err
	}
	response.Message = "注册成功"
	return nil
}

//类型转化
func UserForResponse(userModel *model.User) *user.UserInfoResponse {
	response := &user.UserInfoResponse{}
	response.UserName = userModel.UserName
	response.FirstName = userModel.FirstName
	response.UserId = userModel.ID
	return response
}
