package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

const (
	// Address gRPC服务地址
	Address = "127.0.0.1:8003"
)

func main() {

	// 证书认证-双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, _ := tls.LoadX509KeyPair("two-way/keys2/client.pem", "two-way/keys2/client.key")
	// 创建一个新的、空的 CertPool
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("two-way/keys2/ca.pem")

	//注意这里只能解析pem类型的根证书，所以需要的是ca.pem
	// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	certPool.AppendCertsFromPEM(ca)

	// 构建基于 TLS 的 TransportCredentials 选项
	creds := credentials.NewTLS(&tls.Config{
		// 设置证书链，允许包含一个或多个
		Certificates: []tls.Certificate{cert},
		ServerName:   "ximu.info", //注意这里的参数为配置文件中所允许的ServerName，也就是其中配置的DNS...
		RootCAs:      certPool,
	})

	//1、 建立连接
	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))
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
	fmt.Println("双向认证：调用gRPC方法成功，ProdStock = ", res)

}
