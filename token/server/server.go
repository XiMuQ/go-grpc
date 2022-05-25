package main

import (
	"context"
	"fmt"
	"go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

//1、声明一个server，里面是未实现的字段
type server struct {
	service.UnimplementedQueryUserServer
}

//2、必须要实现在hello.proto里声明的远程调用接口，否则客户端会报：
//rpc error: code = Unimplemented desc = method GetUserInfo not implemente
func (s *server) GetUserInfo(ctx context.Context, in *service.ReqParam) (*service.ResParam, error) {
	return &service.ResParam{Id: in.Id, Name: in.Name, Age: 20, Address: "Beijing"}, nil
}

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:8003"
)

func main() {
	var authInterceptor grpc.UnaryServerInterceptor

	//匿名方法
	authInterceptor = func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		//拦截普通方法请求，验证 Token
		err = Auth(ctx)
		if err != nil {
			return
		}
		// 继续处理请求
		return handler(ctx, req)
	}

	ser := grpc.NewServer(grpc.UnaryInterceptor(authInterceptor))

	//2、注册服务
	service.RegisterQueryUserServer(ser, &server{})

	//3、监听服务端口
	listener, err := net.Listen("tcp", Address)
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	//4、启动服务
	_ = ser.Serve(listener)
}

func Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	var user string
	var password string

	if val, ok := md["user"]; ok {
		user = val[0]
	}
	if val, ok := md["password"]; ok {
		password = val[0]
	}

	if user != "admin" || password != "admin" {
		return status.Errorf(codes.Unauthenticated, "客户端请求的token不合法")
	}
	return nil
}
