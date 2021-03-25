package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RecoveryInterceptor panic时返回Unknown错误吗
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		// p 就是捕获 panic 中的内容
		return status.Errorf(codes.Unknown, "panic 了: %v", p)
	})
}