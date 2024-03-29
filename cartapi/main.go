package main

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	go_micro_service_cart "github.com/caihuahang8/cart/proto"
	"github.com/caihuahang8/cartapi/handler"
	cartApi "github.com/caihuahang8/cartapi/proto"
	"github.com/caihuahang8/common"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/common/log"
	"net"
	"net/http"
)

func main() {
	//注册中心
	reg := etcd.NewRegistry(registry.Addrs("127.0.0.1:2379"))
	//consul := consul2.NewRegistry(func(options *registry.Options) {
	//	options.Addrs = []string{
	//		"127.0.0.1:8500",
	//	}
	//})

	//链路追踪
	t, io, err := common.NewTracer("go.micro.api.cartApi", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//熔断器
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	//启动端口
	go func() {
		err = http.ListenAndServe(net.JoinHostPort("0.0.0.0", " "), hystrixStreamHandler)
		if err != nil {
			log.Error(err)
		}
	}()

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.cartApi"),
		micro.Version("latest"),
		micro.Address("0.0.0.0:8086"),
		//添加 consul 注册中心
		micro.Registry(reg),
		//添加链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//添加熔断
		micro.WrapClient(NewClientHystrixWrapper()),
		//添加负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)

	// Initialise service
	service.Init()

	cartService := go_micro_service_cart.NewCartService("go.micro.service.cart", service.Client())

	//cartService.AddCart(context.TODO(), &go_micro_service_cart.CartInfo{
	//
	//	UserId:    3,
	//	ProductId: 4,
	//	SizeId:    5,
	//	Num:       5,
	//})
	cartService.AddCart(context.TODO(), &go_micro_service_cart.CartInfo{

		UserId:    1,
		ProductId: 2,
		SizeId:    3,
		Num:       4,
	})
	// Register Handler
	if err := cartApi.RegisterCartApiHandler(service.Server(), &handler.CartApi{CartService: cartService}); err != nil {
		log.Error(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	return hystrix.Do(req.Service()+"."+req.Endpoint(), func() error {
		//run 正常执行
		fmt.Println(req.Service() + "." + req.Endpoint())
		return c.Client.Call(ctx, req, rsp, opts...)
	}, func(err error) error {
		fmt.Println(err)
		return err
	})
}

func NewClientHystrixWrapper() client.Wrapper {
	return func(i client.Client) client.Client {
		return &clientWrapper{i}
	}
}
