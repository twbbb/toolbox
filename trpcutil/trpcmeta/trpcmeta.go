// Package trpcmeta trpcmeta
package trpcmeta

import (
	"context"
	"strconv"

	"trpc.group/trpc-go/trpc-go"
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

// GetUid 从context中获取网关透传的uid，没有或不合法时返回0
func GetUid(ctx context.Context) uint64 {
	return getUintFromCtx(ctx, "uid")
}

// GetApp 从context中获取网关透传的app
func GetApp(ctx context.Context) string {
	return getStrFromCtx(ctx, "app")
}

// GetCid 从context中获取网关透传的cid，没有或不合法时返回0
func GetCid(ctx context.Context) uint64 {
	return getUintFromCtx(ctx, "cid")
}

// GetClientIp 从context中获取网关透传的clientIp
func GetClientIp(ctx context.Context) string {
	return getStrFromCtx(ctx, "clientIp")
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

// SetRequestId context中设置requestId
func SetRequestId(ctx context.Context, sthid string) {
	trpc.SetMetaData(ctx, "requestId", []byte(sthid))
}

// GetCityCode context 中获取 citycode
func GetCityCode(ctx context.Context) string {
	return CommonGet(ctx, []string{"cityCode", "CityCode", "citycode", "Citycode"})
}

// CommonGet 通用获取某个 key 的值
func CommonGet(ctx context.Context, keys []string) string {
	for _, key := range keys {
		v := getStrFromCtx(ctx, key)
		if v != "" {
			return v
		}
	}
	return ""
}

func getStrFromCtx(ctx context.Context, label string) string {
	if value := trpc.GetMetaData(ctx, label); value != nil {
		return string(value)
	}
	return ""
}

func getUintFromCtx(ctx context.Context, label string) uint64 {
	v, err := strconv.ParseUint(getStrFromCtx(ctx, label), 10, 64)
	if err != nil {
		return 0
	}
	return v
}
