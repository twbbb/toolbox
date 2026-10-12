module github.com/twbbb/toolbox/bus/ocr

go 1.23.4

require trpc.group/trpc-go/trpc-go v1.1.0

require (
	github.com/golang/snappy v0.0.3 // indirect
	github.com/google/flatbuffers v2.0.0+incompatible // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/lestrrat-go/strftime v1.0.6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180228061459-e0a39a4cb421 // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	trpc.group/trpc/trpc-protocol/pb/go/trpc v1.0.0 // indirect
)

// fasthttp 钉在 v1.43.0：公网 trpc-go v1.1.0 按这个版本写的 RoundTripper。
replace github.com/valyala/fasthttp => github.com/valyala/fasthttp v1.43.0
