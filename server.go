package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/xhyonline/grpc-middleware/auth"
	"github.com/xhyonline/grpc-middleware/pb"
	"github.com/xhyonline/grpc-middleware/recovery"
	zap_log "github.com/xhyonline/grpc-middleware/zap-log"
	"github.com/xhyonline/xutil/helper"
	"google.golang.org/grpc"
	"net"
)

type SimpleServer struct {
}

func (s *SimpleServer) Hello(ctx context.Context, empty *pb.Empty) (*pb.Response, error) {
	// 获取 10 以内随机数,如果大于 5 直接报错
	// 用于测试 panic 拦截器
	if helper.GetRandom(10)>5{
		panic("触发了 panic 请观察 grpc 服务是否中断")
	}
	fmt.Println("成功请求到了 Hello ")
	// 此时你就能获取鉴权中设置的 key
	fmt.Println(ctx.Value("1"))
	return &pb.Response{}, nil
}


func main() {
	g := grpc.NewServer(grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
		// 添加鉴权拦截器
		grpc_auth.UnaryServerInterceptor(auth.AuthInterceptor),
		// 添加日志分析拦截器
		grpc_zap.UnaryServerInterceptor(zap_log.ZapInterceptor()),
		// 添加 panic 拦截器
		grpc_recovery.UnaryServerInterceptor(recovery.RecoveryInterceptor()),
	)))
	pb.RegisterSimpleServer(g, new(SimpleServer))
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	err = g.Serve(l)
	if err != nil {
		panic(err)
	}
}
