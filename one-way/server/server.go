package main

import (
	"context"
	"go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
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

func main() {

	//配置 TLS认证相关文件
	creds, err := credentials.NewServerTLSFromFile("one-way/keys/server.pem", "one-way/keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}

	//1、创建服务，并开启TLS认证
	//ser := grpc.NewServer()
	ser := grpc.NewServer(grpc.Creds(creds))

	//2、注册服务
	service.RegisterQueryUserServer(ser, &server{})

	//3、监听服务端口
	listener, err := net.Listen("tcp", ":8002")
	if err != nil {
		log.Fatal("服务监听端口失败", err)
	}

	//4、启动服务
	_ = ser.Serve(listener)
}
