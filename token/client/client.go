package main

import (
	"context"
	"fmt"
	"go-grpc/service"
	"go-grpc/token/service/impl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:8003"
)

func main() {

	user := &impl.Authentication{
		User:     "admin",
		Password: "admin2",
	}
	//1、 建立连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithPerRPCCredentials(user))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	request := &service.ReqParam{
		Id:   123,
		Name: "西木",
	}

	// 2. 调用 hello_grpc.pb.go 中的NewQueryUserClient方法建立客户端
	query := service.NewQueryUserClient(conn)

	//3、调用rpc方法
	res, err := query.GetUserInfo(context.Background(), request)

	if err != nil {
		log.Fatal("调用gRPC方法错误: ", err)
	}
	fmt.Println("Token：调用gRPC方法成功，ProdStock = ", res)

}
