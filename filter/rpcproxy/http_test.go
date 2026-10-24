package rpcproxy

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// fakeServer 回显收到的 body，记调用次数
func fakeServer(t *testing.T, calls *int) *httptest.Server {
	t.Helper()
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*calls++
		body, _ := io.ReadAll(r.Body)
		if r.URL.Path == "/text" {
			w.Header().Set("Content-Type", "text/plain")
			_, _ = w.Write([]byte("plain:" + string(body)))
			return
		}
		if r.URL.Query().Get("access_token") == "old" {
			w.WriteHeader(http.StatusForbidden)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"got":` + string(body) + `,"token":"` + r.URL.Query().Get("access_token") + `"}`))
	}))
	t.Cleanup(s.Close)
	return s
}

func post(t *testing.T, client *http.Client, u, body string) (int, string) {
	t.Helper()
	rsp, err := client.Post(u, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer rsp.Body.Close()
	data, _ := io.ReadAll(rsp.Body)
	return rsp.StatusCode, string(data)
}

func TestHTTPCacheHit(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: true, LocalFilePath: dir})
	calls := 0
	s := fakeServer(t, &calls)
	client := &http.Client{Transport: Transport(nil)}

	code, body := post(t, client, s.URL+"/cgi-bin/msg?access_token=new&nonce=1", `{"a":1,"nonce":"x"}`)
	if calls != 1 || code != 200 || body != `{"got":{"a":1,"nonce":"x"},"token":"new"}` {
		t.Fatalf("first call should reach server with body intact, calls=%d code=%d body=%s", calls, code, body)
	}
	// query 和 body 里的 nonce 变化不影响命中
	code, body = post(t, client, s.URL+"/cgi-bin/msg?nonce=2&access_token=new", `{"a":1,"nonce":"y"}`)
	if calls != 1 || code != 200 || body != `{"got":{"a":1,"nonce":"x"},"token":"new"}` {
		t.Fatalf("want cache hit, calls=%d code=%d body=%s", calls, code, body)
	}
	// access_token 不同是另一份文件；非 2xx 也照常缓存
	for i := 0; i < 2; i++ {
		code, _ = post(t, client, s.URL+"/cgi-bin/msg?access_token=old", `{"a":1}`)
		if calls != 2 || code != http.StatusForbidden {
			t.Fatalf("round %d: different access_token should be its own file, calls=%d code=%d", i, calls, code)
		}
	}
	files, _ := filepath.Glob(filepath.Join(dir, "http.127.0.0.1*.cgi-bin_msg_*.json"))
	if len(files) != 2 {
		t.Fatalf("want 2 cache files named http.<host>.cgi-bin_msg_<md5>.json, got %v", files)
	}
}

func TestHTTPEditedRsp(t *testing.T) {
	dir := t.TempDir()
	withConfig(t, Config{Enable: true, LocalFilePath: dir})
	calls := 0
	s := fakeServer(t, &calls)
	client := &http.Client{Transport: Transport(nil)}

	post(t, client, s.URL+"/text", "hello")
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Fatalf("want one cache file, got %d", len(files))
	}
	file := filepath.Join(dir, files[0].Name())
	data, _ := os.ReadFile(file)
	if !strings.Contains(string(data), `"bodyText": "plain:hello"`) {
		t.Fatalf("non-JSON body should be stored as bodyText: %s", data)
	}
	// 手写替身：只写 rsp，status 和 JSON body
	if err := os.WriteFile(file, []byte(`{"req":{},"rsp":{"status":500,"body":{"errcode":40001}}}`), 0644); err != nil {
		t.Fatal(err)
	}
	code, body := post(t, client, s.URL+"/text", "hello")
	if calls != 1 || code != 500 || body != `{"errcode":40001}` {
		t.Fatalf("want edited rsp, calls=%d code=%d body=%s", calls, code, body)
	}
	// rsp 不是对象时回退到真实调用
	if err := os.WriteFile(file, []byte(`{"req":{},"rsp":"bad"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if code, body = post(t, client, s.URL+"/text", "hello"); calls != 2 || body != "plain:hello" {
		t.Fatalf("invalid rsp should fall through, calls=%d code=%d body=%s", calls, code, body)
	}
}

func TestHTTPDisabledAndNoCache(t *testing.T) {
	dir := t.TempDir()
	calls := 0
	s := fakeServer(t, &calls)
	client := &http.Client{Transport: Transport(nil)}

	withConfig(t, Config{Enable: false, LocalFilePath: dir})
	post(t, client, s.URL+"/a/b", `{}`)
	post(t, client, s.URL+"/a/b", `{}`)
	if calls != 2 {
		t.Fatalf("disabled should always call server, calls=%d", calls)
	}

	host := strings.TrimPrefix(s.URL, "http://")
	withConfig(t, Config{Enable: true, LocalFilePath: dir, NoCacheMethods: []string{"http." + host + "/a_b"}})
	post(t, client, s.URL+"/a/b", `{}`)
	post(t, client, s.URL+"/a/b", `{}`)
	if files, _ := os.ReadDir(dir); calls != 4 || len(files) != 0 {
		t.Fatalf("no-cache method should always call server and not write, calls=%d files=%d", calls, len(files))
	}
}
