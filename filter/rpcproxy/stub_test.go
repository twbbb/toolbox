package rpcproxy

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestHTTPStubHit(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: true, LocalFilePath: dir, ReadOnly: true})
	calls := 0
	s := fakeServer(t, &calls)
	client := &http.Client{Transport: Transport(nil)}

	u := s.URL + "/wxa/msg_sec_check?access_token=old&nonce=9"
	file, err := HTTPStubFile(dir, "post", u, []byte(`{"content":"hi"}`))
	if err != nil {
		t.Fatal(err)
	}
	rsp := &HTTPRsp{Body: json.RawMessage(`{"errcode":40001}`)}
	if err := WriteStub(file, &HTTPReq{Method: "POST", URL: u}, rsp); err != nil {
		t.Fatal(err)
	}
	// nonce 不同照样命中
	code, body := post(t, client, strings.Replace(u, "nonce=9", "nonce=1", 1), `{"content":"hi"}`)
	if calls != 0 || code != 200 || body != `{"errcode":40001}` {
		t.Fatalf("want stub hit, calls=%d code=%d body=%s", calls, code, body)
	}
	// 没有替身的走真实调用，readOnly 不写缓存
	post(t, client, s.URL+"/wxa/msg_sec_check?access_token=new", `{"content":"hi"}`)
	if files, _ := os.ReadDir(dir); calls != 1 || len(files) != 1 {
		t.Fatalf("miss should call server and not write, calls=%d files=%d", calls, len(files))
	}
}

func TestRPCStubHit(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: true, LocalFilePath: dir, ReadOnly: true})
	const svc = "trpc.a.b.S"
	file, err := StubFile(dir, svc, "Get", &funReq{UserID: "u1", Nonce: "whatever"}, map[string]string{"sthid": "@u9"})
	if err != nil {
		t.Fatal(err)
	}
	if err := WriteStub(file, &funReq{UserID: "u1"}, &funRsp{Name: "stub"}); err != nil {
		t.Fatal(err)
	}
	calls := 0
	if rsp := call(t, svc, "Get", "@u9", "n1", &calls, nil); calls != 0 || rsp.Name != "stub" {
		t.Fatalf("want stub hit, calls=%d rsp=%+v", calls, rsp)
	}
	// 别的用户没有替身，真实调用，不写缓存
	call(t, svc, "Get", "@u1", "n1", &calls, nil)
	if files, _ := os.ReadDir(dir); calls != 1 || len(files) != 1 {
		t.Fatalf("miss should call next and not write, calls=%d files=%d", calls, len(files))
	}
}
