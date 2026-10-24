package rpcproxy

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"trpc.group/trpc-go/trpc-go/codec"
)

type funReq struct {
	UserID string `json:"user_id"`
	Nonce  string `json:"nonce"`
}

type funRsp struct {
	Name string `json:"name"`
}

func withConfig(t *testing.T, c Config) {
	t.Helper()
	SetConfig(c)
	t.Cleanup(func() { SetConfig(Config{}) })
}

func call(t *testing.T, service, method, sthid, nonce string, nextCalls *int, nextErr error) *funRsp {
	t.Helper()
	ctx, msg := codec.WithNewMessage(context.Background())
	msg.WithCalleeServiceName(service)
	msg.WithClientRPCName("/" + service + "/" + method)
	msg.WithClientMetaData(codec.MetaData{
		"sthid":     []byte(sthid),
		"requestId": []byte(fmt.Sprint(*nextCalls)),
		"bin":       []byte{0xff, 0xfe},
	})
	req := &funReq{UserID: "u1", Nonce: nonce}
	rsp := &funRsp{}
	err := clientFilter(ctx, req, rsp, func(_ context.Context, _, r interface{}) error {
		*nextCalls++
		if nextErr != nil {
			return nextErr
		}
		r.(*funRsp).Name = sthid
		return nil
	})
	if nextErr != nil {
		if err == nil {
			t.Fatal("want next error")
		}
		return rsp
	}
	if err != nil {
		t.Fatal(err)
	}
	return rsp
}

func TestCacheHitAndForce(t *testing.T) {
	arg := Config{Enable: true, LocalFilePath: t.TempDir()}
	withConfig(t, arg)
	calls := 0
	const svc = "trpc.a.b.S"

	rsp := call(t, svc, "Get", "@u1", "n1", &calls, nil)
	if calls != 1 || rsp.Name != "@u1" {
		t.Fatalf("first call should go downstream, calls=%d rsp=%+v", calls, rsp)
	}
	// nonce、requestId 变化不影响命中
	rsp = call(t, svc, "Get", "@u1", "n2", &calls, nil)
	if calls != 1 || rsp.Name != "@u1" {
		t.Fatalf("want cache hit, calls=%d rsp=%+v", calls, rsp)
	}
	rsp = call(t, svc, "Get", "@u2", "n2", &calls, nil)
	if calls != 2 || rsp.Name != "@u2" {
		t.Fatalf("different sthid should miss cache, calls=%d rsp=%+v", calls, rsp)
	}

	arg.Force = true
	withConfig(t, arg)
	call(t, svc, "Get", "@u1", "n2", &calls, nil)
	if calls != 3 {
		t.Fatalf("force should skip cache, calls=%d", calls)
	}
}

func TestErrorNotCached(t *testing.T) {
	withConfig(t, Config{Enable: true, LocalFilePath: t.TempDir()})
	calls := 0
	call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, errors.New("down"))
	call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, errors.New("down"))
	if calls != 2 {
		t.Fatalf("error should not be cached, calls=%d", calls)
	}
}

func TestNoCacheMethods(t *testing.T) {
	const svc = "trpc.a.b.S"
	cases := []struct {
		name    string
		methods []string
		noCache bool
	}{
		{"method", []string{"Get"}, true},
		{"service/method", []string{svc + "/Get"}, true},
		{"other service", []string{"trpc.other.S/Get"}, false},
		{"other method", []string{"List"}, false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			withConfig(t, Config{Enable: true, LocalFilePath: dir, NoCacheMethods: c.methods})
			calls := 0
			call(t, svc, "Get", "@u1", "n1", &calls, nil)
			call(t, svc, "Get", "@u1", "n1", &calls, nil)
			files, _ := os.ReadDir(dir)
			if c.noCache && (calls != 2 || len(files) != 0) {
				t.Fatalf("no-cache method should always call downstream and not write cache, calls=%d files=%d",
					calls, len(files))
			}
			if !c.noCache && (calls != 1 || len(files) != 1) {
				t.Fatalf("method should be cached, calls=%d files=%d", calls, len(files))
			}
		})
	}
}

func TestDisabledAndSkipService(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: false, LocalFilePath: dir})
	calls := 0
	call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, nil)
	call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, nil)
	if calls != 2 {
		t.Fatalf("disabled filter should always call downstream, calls=%d", calls)
	}

	withConfig(t, Config{Enable: true, LocalFilePath: dir})
	call(t, "trpc.mysql.tencent_health.th_jingqi_server.jingqi", "Get", "@u1", "n1", &calls, nil)
	call(t, "trpc.mysql.tencent_health.th_jingqi_server.jingqi", "Get", "@u1", "n1", &calls, nil)
	if calls != 4 {
		t.Fatalf("storage client should skip cache, calls=%d", calls)
	}
	if !isSkipService("trpc.mysql.tencent_health.th_jingqi_server.jingqi") || isSkipService("trpc.health.todo.HealthTodo") {
		t.Fatal("skip rule wrong")
	}
}

func TestInvalidCacheFallsThrough(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: true, LocalFilePath: dir})
	calls := 0
	call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, nil)
	files, err := os.ReadDir(dir)
	if err != nil || len(files) != 1 {
		t.Fatalf("want one cache file, err=%v files=%d", err, len(files))
	}
	path := dir + "/" + files[0].Name()
	if err := os.WriteFile(path, []byte(`{"req":{},"rsp":"bad"}`), 0644); err != nil {
		t.Fatal(err)
	}
	rsp := call(t, "trpc.a.b.S", "Get", "@u1", "n1", &calls, nil)
	if calls != 2 || rsp.Name != "@u1" {
		t.Fatalf("invalid rsp should fall through, calls=%d rsp=%+v", calls, rsp)
	}
}
