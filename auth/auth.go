package auth

import (
	"context"
	"errors"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TokenInfo 用户信息
type TokenInfo struct {
	ID    string
	Roles []string
}

// AuthInterceptor 认证拦截器，对以authorization为头部，形式为`bearer token`的Token进行验证
func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	// 解析认证 token
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "token 认证失败 %v", err)
	}
	//使用context.WithValue添加了值后，可以用Value(key)方法获取值
	// 因此,通过鉴权后,你可以根据 Hello 方法的上下文获取
	newCtx := context.WithValue(ctx, tokenInfo.ID, tokenInfo)
	return newCtx, nil
}

//解析token，并进行验证
func parseToken(token string) (TokenInfo, error) {
	var tokenInfo TokenInfo
	if token == "grpc.auth.token" {
		tokenInfo.ID = "1"
		tokenInfo.Roles = []string{"admin"}
		return tokenInfo, nil
	}
	return tokenInfo, errors.New("Token无效: bearer " + token)
}

//从token中获取用户唯一标识
func userClaimFromToken(tokenInfo TokenInfo) string {
	return tokenInfo.ID
}
