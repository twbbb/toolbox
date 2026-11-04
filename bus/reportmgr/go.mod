module github.com/twbbb/toolbox/bus/reportmgr

go 1.23

// fasthttp 钉在 v1.43.0：公网 trpc-go v1.1.0 按这个版本写的 RoundTripper。
replace github.com/valyala/fasthttp => github.com/valyala/fasthttp v1.43.0

replace github.com/twbbb/toolbox/component/aesutil => ../../component/aesutil
