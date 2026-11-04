package rpcproxy

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"trpc.group/trpc-go/trpc-go/log"
)

// httpServicePrefix HTTP 下游的服务名前缀，服务名是 http.<host>
const httpServicePrefix = "http."

// queryIgnoreKeys 每次都会变的 query 参数，不参与缓存 key
var queryIgnoreKeys = []string{"nonce", "timestamp", "signature"}

// Transport 给 http.Client 的下游调用加同样的文件缓存，用法：
//
//	client := &http.Client{Timeout: 5 * time.Second, Transport: rpcproxy.Transport(nil)}
//
// base 为 nil 时用 http.DefaultTransport。DOCKER_ENV 非空（容器环境）时直接返回 base，线上不生效。
func Transport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	if os.Getenv("DOCKER_ENV") != "" {
		return base
	}
	return &transport{base: base}
}

type transport struct {
	base http.RoundTripper
}

// HTTPReq HTTP 下游缓存文件里的 req，只用来查看
type HTTPReq struct {
	Method   string          `json:"method"`
	URL      string          `json:"url"`
	Body     json.RawMessage `json:"body,omitempty"`     // 请求体是 JSON 时原样放这里
	BodyText string          `json:"bodyText,omitempty"` // 不是 JSON 时放这里
}

// HTTPRsp HTTP 下游缓存文件里的 rsp，命中时按它返回；手写替身时 Status 缺省为 200
type HTTPRsp struct {
	Status   int               `json:"status,omitempty"` // 缺省 200
	Header   map[string]string `json:"header,omitempty"`
	Body     json.RawMessage   `json:"body,omitempty"`     // 响应体是 JSON 时原样放这里
	BodyText string            `json:"bodyText,omitempty"` // 不是 JSON 时放这里
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	c, _ := conf.Load().(*Config)
	if c == nil || !c.Enable {
		return t.base.RoundTrip(req)
	}
	arg := *c
	ctx := req.Context()
	serviceName := httpServicePrefix + req.URL.Host
	method := httpMethodName(req.URL.Path)

	reqBody, err := readBody(&req.Body)
	if err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|read req body fail, skip cache|%s|%s|%v", serviceName, method, err)
		return t.base.RoundTrip(req)
	}
	file := cacheFile(ctx, &arg, serviceName, method, httpCacheKey(req, reqBody))
	if data, ok := readCache(ctx, &arg, file); ok {
		if rsp, err := buildResponse(req, data); err != nil {
			log.ErrorContextf(ctx, "rpcproxy cache|rsp unmarshal fail, ignore|%s|%v", file, err)
		} else {
			return rsp, nil
		}
	}
	rsp, err := t.base.RoundTrip(req)
	if err != nil || file == "" {
		return rsp, err
	}
	rspBody, err := readBody(&rsp.Body)
	if err != nil {
		log.ErrorContextf(ctx, "rpcproxy cache|read rsp body fail, skip cache|%s|%v", file, err)
		return rsp, nil
	}
	reqJSON, _ := json.Marshal(newHTTPCacheReq(req, reqBody))
	rspJSON, _ := json.Marshal(newHTTPCacheRsp(rsp, rspBody))
	writeCache(ctx, &arg, file, reqJSON, rspJSON)
	return rsp, nil
}

// httpMethodName path 去掉首尾 / 后把 / 换成 _，如 /cgi-bin/stable_token → cgi-bin_stable_token
func httpMethodName(p string) string {
	p = strings.Trim(p, "/")
	if p == "" {
		return "_"
	}
	return strings.ReplaceAll(p, "/", "_")
}

// readBody 读出整个 body 再放回去，后面的调用还能再读
func readBody(body *io.ReadCloser) ([]byte, error) {
	if *body == nil || *body == http.NoBody {
		return nil, nil
	}
	data, err := io.ReadAll(*body)
	_ = (*body).Close()
	*body = io.NopCloser(bytes.NewReader(data))
	return data, err
}

// httpCacheKey method、host、path、去掉易变参数并排好序的 query、去掉易变字段的 body 一起取 md5；
// access_token 这类参数参与 key，旧 token 和新 token 命中不同的文件
func httpCacheKey(req *http.Request, body []byte) string {
	q := req.URL.Query()
	for _, k := range queryIgnoreKeys {
		q.Del(k)
	}
	s := req.Method + " " + req.URL.Host + req.URL.Path + "?" + q.Encode() + "\n" + stripIgnored(string(body))
	return fmt.Sprintf("%x", md5.Sum([]byte(s)))
}

func newHTTPCacheReq(req *http.Request, body []byte) *HTTPReq {
	r := &HTTPReq{Method: req.Method, URL: req.URL.String()}
	if json.Valid(body) {
		r.Body = body
	} else {
		r.BodyText = string(body)
	}
	return r
}

func newHTTPCacheRsp(rsp *http.Response, body []byte) *HTTPRsp {
	r := &HTTPRsp{Status: rsp.StatusCode, Header: map[string]string{}}
	for k := range rsp.Header {
		if k != "Content-Length" {
			r.Header[k] = rsp.Header.Get(k)
		}
	}
	if json.Valid(body) {
		r.Body = body
	} else {
		r.BodyText = string(body)
	}
	return r
}

// buildResponse 按缓存的 rsp 拼出响应；status 缺省当 200
func buildResponse(req *http.Request, data []byte) (*http.Response, error) {
	c := &HTTPRsp{}
	if err := json.Unmarshal(data, c); err != nil {
		return nil, err
	}
	if c.Status == 0 {
		c.Status = http.StatusOK
	}
	body := []byte(c.BodyText)
	if len(c.Body) > 0 {
		// 缓存文件是缩进写的，返回前压回紧凑格式
		buf := &bytes.Buffer{}
		if err := json.Compact(buf, c.Body); err != nil {
			return nil, err
		}
		body = buf.Bytes()
	}
	h := http.Header{}
	for k, v := range c.Header {
		h.Set(k, v)
	}
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", c.Status, http.StatusText(c.Status)),
		StatusCode:    c.Status,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        h,
		Body:          io.NopCloser(bytes.NewReader(body)),
		ContentLength: int64(len(body)),
		Request:       req,
	}, nil
}
