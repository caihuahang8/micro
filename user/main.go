package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"user/domain/repository"
	"user/domain/service"
	"user/handler"
	user "user/proto"
)

func main() {
	// 2.注册中心
	consul := consul2.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"localhost:8500",
		}
	})
	//服务参数设置
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address(":9080"),
		//添加consul 注册中心
		micro.Registry(consul),
	)
	//初始化服务
	srv.Init()

	//创建数据库连接
	db, err := gorm.Open("mysql", "root:123456@/micro?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	db.SingularTable(true)

	//只执行一次，数据表初始化
	//rp:=repository.NewUserRepository(db)
	//rp.InitTable()

	//创建服务实例
	userDataService := service.NewUserDataService(repository.NewUserRepository(db))
	//注册Handler
	err = user.RegisterUserHandler(srv.Server(), &handler.User{UserService: userDataService})
	if err != nil {
		fmt.Println(err)
	}

	// Run service
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
