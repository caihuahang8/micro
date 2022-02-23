package handler

import (
	"context"
	"encoding/json"
	"errors"
	cart "github.com/caihuahang8/cart/proto"
	cartApi "github.com/caihuahang8/cartapi/proto"
	"github.com/prometheus/common/log"
	"strconv"
)

type CartApi struct {
	CartService cart.CartService
}

// @Router /users/{id} [delete]
// CartApi.Call 通过API向外暴露为/cartApi/findAll，接收http请求
func (e *CartApi) FindAll(ctx context.Context, req *cartApi.Request, rsp *cartApi.Response) error {
	log.Info("接受到 /cartApi/findAll 访问请求")
	userId := req.UserId
	if _, ok := req.Get["user_id"]; !ok {
		return errors.New("参数异常")
	}
	userIdString := req.Get["user_id"].Values[0]
	userId, err := strconv.ParseInt(userIdString, 10, 64)
	if err != nil {
		return err
	}
	//获取购物车所有商品
	cartAll, err := e.CartService.GetAll(context.TODO(), &cart.CartFindAll{UserId: userId})
	//数据类型转化
	b, err := json.Marshal(cartAll)
	if err != nil {
		return err
	}
	rsp.StatusCode = 200
	rsp.Body = string(b)
	return nil
}
