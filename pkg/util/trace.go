package util

import (
	"fmt"
	"github.com/cloudwego/kitex/client"
	server "github.com/cloudwego/kitex/server"
	internal_opentracing "github.com/kitex-contrib/tracer-opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"io"
)

func Trace(service string) (client.Suite, io.Closer) {
	closer := build(service)
	return internal_opentracing.NewDefaultClientSuite(), closer
}

func SrvTrace(service string) (server.Suite, io.Closer) {
	closer := build(service)
	return internal_opentracing.NewDefaultServerSuite(), closer
}
func build(service string) io.Closer {
	// 可以配置成 我们 变量, 但是win十分.....
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		panic(err)
	}
	//cfg := &jaegercfg.Configuration{
	//	ServiceName: service,
	//	Disabled:    false,
	//	RPCMetrics:  false,
	//	Gen128Bit:   false,
	//	Tags:        nil,
	//	Sampler: &jaegercfg.SamplerConfig{
	//		Type:                     "const",
	//		Param:                    1,
	//		SamplingServerURL:        "",
	//		SamplingRefreshInterval:  0,
	//		MaxOperations:            0,
	//		OperationNameLateBinding: false,
	//		Options:                  nil,
	//	},
	//	Reporter: &jaegercfg.ReporterConfig{
	//		LogSpans: true,
	//	},
	//	Headers:             nil,
	//	BaggageRestrictions: nil,
	//	Throttler:           nil,
	//}

	tracer, closer, err := cfg.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger: %v\n", err))
	}
	opentracing.InitGlobalTracer(tracer)
	return closer
}
