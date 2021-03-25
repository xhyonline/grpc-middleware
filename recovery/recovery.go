package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

// RecoveryInterceptor panic时返回Unknown错误吗
func RecoveryInterceptor() grpc_recovery.Option {
	return grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
		buf := make([]byte, 4096)
		// 抛出服务端 调用栈的轨迹
		// runtime.stack 详情参考资料
		// https://colobu.com/2016/12/21/how-to-dump-goroutine-stack-traces/
		num := runtime.Stack(buf, false)
		// p 就是捕获 panic 中的内容
		return status.Errorf(codes.Unknown, "panic 了: %v %s", p,string(buf[:num]))
	})
}