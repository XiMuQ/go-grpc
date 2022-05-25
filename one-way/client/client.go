package main

import (
	"context"
	"fmt"
	"go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"log"
)

func main() {

	// TLS连接  注意serverNameOverride填写的是在server.conf的[alt_names]里的DNS.X地址
	creds, err := credentials.NewClientTLSFromFile("one-way/keys/server.pem", "localhost")
	if err != nil {
		grpclog.Fatalf("Failed to create TLS credentials %v", err)
	}

	//1、 建立连接
	//conn, err := grpc.Dial(":8002", grpc.WithInsecure())
	conn, err := grpc.Dial(":8002", grpc.WithTransportCredentials(creds))
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
	fmt.Println("单向认证：调用gRPC方法成功，ProdStock = ", res)

}
