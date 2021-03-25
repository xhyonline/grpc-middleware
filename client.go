package main

import (
	"fmt"
	"github.com/xhyonline/grpc-middleware/pb"
	"google.golang.org/grpc"
	"context"
)

// Token token认证
type Token struct {
	Value string
}

// 以下两个方法均是为了实现 grpc.WithPerRPCCredentials 入参的接口

// GetRequestMetadata 获取当前请求认证所需的元数据 (实现接口1)
func (t *Token) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{"authorization": t.Value}, nil
}

// RequireTransportSecurity 是否需要基于 TLS 认证进行安全传输 (实现接口2)
// 这里我就不用 tls 了,否则还需要通过下面去生成证书
// creds, err := credentials.NewClientTLSFromFile("../pkg/tls/server.pem", "go-grpc-example")
// grpc.WithTransportCredentials(creds)
func (t *Token) RequireTransportSecurity() bool {
	// 不需要证书
	return false
}

func main() {
	token:=&Token{Value: "bearer grpc.auth.token"}
	connect, err := grpc.Dial("127.0.0.1:8080",grpc.WithInsecure(), grpc.WithPerRPCCredentials(token))
	if err != nil {
		panic(err)
	}
	defer connect.Close()
	client:=pb.NewSimpleClient(connect)
	result,err:=client.Hello(context.Background(),&pb.Empty{})
	if err!=nil{
		panic(err)
	}
	fmt.Println(result)
}