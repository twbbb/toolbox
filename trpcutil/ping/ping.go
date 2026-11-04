// Package ping 网关探活用的默认 Ping 实现
package ping

import "context"

// Default 默认 Ping，嵌进业务 service 即可：
//
//	type svc struct{ ping.Default[pb.PingReq, pb.PingRsp] }
//
// 只说明进程活着，不查 DB、Redis。
type Default[Req, Rsp any] struct{}

// Ping 返回空响应
func (Default[Req, Rsp]) Ping(_ context.Context, _ *Req) (*Rsp, error) {
	return new(Rsp), nil
}
