package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"go-grpc/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
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

	// TLS认证
	//从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, _ := tls.LoadX509KeyPair("two-way/keys2/server.pem", "two-way/keys2/server.key")
	//初始化一个CertPool
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("two-way/keys2/ca.pem")

	//注意这里只能解析pem类型的根证书，所以需要的是ca.pem
	//解析传入的证书，解析成功会将其加到池子中
	certPool.AppendCertsFromPEM(ca)

	//构建基于TLS的TransportCredentials选项
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},        //服务端证书链，可以有多个
		ClientAuth:   tls.RequireAndVerifyClientCert, //要求必须验证客户端证书
		ClientCAs:    certPool,                       //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})

	//1、创建服务，并开启TLS认证
	//ser := grpc.NewServer()
	//ser := grpc.NewServer(grpc.Creds(creds), grpc.KeepaliveParams(keepalive.ServerParameters{
	//	MaxConnectionIdle: 5 * time.Minute, //这个连接最大的空闲时间，超过就释放，解决proxy等到网络问题（不通知grpc的client和server）
	//}))
	ser := grpc.NewServer(grpc.Creds(creds))

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
