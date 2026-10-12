module github.com/twbbb/toolbox/errcode

go 1.20

require trpc.group/trpc-go/trpc-go v1.1.0

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	trpc.group/trpc/trpc-protocol/pb/go/trpc v1.0.0 // indirect
)

// fasthttp 钉在 v1.43.0：公网 trpc-go v1.1.0 按这个版本写的 RoundTripper。
replace github.com/valyala/fasthttp => github.com/valyala/fasthttp v1.43.0
