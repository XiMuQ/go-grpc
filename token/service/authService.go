package service

import "context"

type PerRPCCredentials interface {

	//GetRequestMetadata 方法返回认证需要的必要信息
	GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error)

	//RequireTransportSecurity 方法表示是否启用安全链接，在生产环境中，一般都是启用的，但为了测试方便，暂时这里不启用
	RequireTransportSecurity() bool
}
