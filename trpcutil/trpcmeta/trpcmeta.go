// Package trpcmeta trpcmeta
package trpcmeta

import (
	"context"

	"git.code.oa.com/trpc-go/trpc-go"
)

// GetOpenID 从context中获取openID
func GetOpenID(ctx context.Context) string {
	return getStrFromCtx(ctx, "openId")
}

// GetUserID 从context中获取userID
func GetUserID(ctx context.Context) string {
	return getStrFromCtx(ctx, "userId")
}

// GetHCUnionId 从context中获取HCUnionId
func GetHCUnionId(ctx context.Context) string {
	return getStrFromCtx(ctx, "HCUnionId")
}

// GetSoid 从context中获取soid
func GetSoid(ctx context.Context) string {
	return getStrFromCtx(ctx, "soid")
}

// GetSthid 从context中获取sthid
func GetSthid(ctx context.Context) string {
	return getStrFromCtx(ctx, "sthid")
}

// GetRequestId 从context中获取requestId
func GetRequestId(ctx context.Context) string {
	return getStrFromCtx(ctx, "requestId")
}

// GetSubBusinessid 从context中获取GetSubBusinessid
func GetSubBusinessid(ctx context.Context) string {
	return getStrFromCtx(ctx, "sub-businessid")
}

// GetVersion 从context中获取version
func GetVersion(ctx context.Context) string {
	return getStrFromCtx(ctx, "version")
}

// GetScene 从context中获取GetScene
func GetScene(ctx context.Context) string {
	return getStrFromCtx(ctx, "scene")
}

// GetChannel 从context中获取GetChannel
func GetChannel(ctx context.Context) string {
	return getStrFromCtx(ctx, "channel")
}

// SetSthid context中设置Sthid
func SetSthid(ctx context.Context, sthid string) {
	trpc.SetMetaData(ctx, "sthid", []byte(sthid))
}
func getStrFromCtx(ctx context.Context, label string) string {
	if value := trpc.GetMetaData(ctx, label); value != nil {
		return string(value)
	}
	return ""
}
