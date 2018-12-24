package kernal

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	// "github.com/opentracing/opentracing-go/log"
	jaeger "github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

//MDReaderWriter metadata Reader and Writer
type MDReaderWriter struct {
	metadata.MD
}

// ForeachKey implements ForeachKey of opentracing.TextMapReader
func (c MDReaderWriter) ForeachKey(handler func(key, val string) error) error {
	for k, vs := range c.MD {
		for _, v := range vs {
			if err := handler(k, v); err != nil {
				return err
			}
		}
	}
	return nil
}

// Set implements Set() of opentracing.TextMapWriter
func (c MDReaderWriter) Set(key, val string) {
	key = strings.ToLower(key)
	c.MD[key] = append(c.MD[key], val)
}

// NewJaegerTracer NewJaegerTracer for current service
func NewJaegerTracer(serviceName string, jagentHost string) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jcfg := jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  jagentHost,
		},
	}

	tracer, closer, err = jcfg.New(
		serviceName,
		jaegercfg.Logger(jaeger.StdLogger),
	)
	if err != nil {
		return
	}

	opentracing.SetGlobalTracer(tracer)
	return
}

func jaegerServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	tracer opentracing.Tracer) context.Context {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	spanContext, err := tracer.Extract(opentracing.TextMap, MDReaderWriter{md})
	if err != nil && err != opentracing.ErrSpanContextNotFound {
		grpclog.Errorf("extract from metadata err: %v", err)
		fmt.Println("[tzj] extract from metadata err: ", err)
	} else {
		fmt.Println("[tzj] extract from metadata : ", spanContext)
		span := tracer.StartSpan(
			info.FullMethod,
			ext.RPCServerOption(spanContext),
			opentracing.Tag{Key: string(ext.Component), Value: "gRPC"},
			ext.SpanKindRPCServer,
		)
		defer span.Finish()

		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	return ctx
}

func jaegerClientInterceptor(
	ctx context.Context,
	method string,
	req, resp interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	tracer opentracing.Tracer) context.Context {

	var parentCtx opentracing.SpanContext
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx = parent.Context()
	}

	clientSpan := tracer.StartSpan(
		method,
		opentracing.ChildOf(parentCtx),
		ext.SpanKindRPCClient,
		// gRPCComponentTag,
	)
	defer clientSpan.Finish()
	return injectSpanContext(ctx, tracer, clientSpan)
}

func injectSpanContext(ctx context.Context, tracer opentracing.Tracer, clientSpan opentracing.Span) context.Context {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}
	mdWriter := MDReaderWriter{md}
	err := tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, mdWriter)
	// We have no better place to record an error than the Span itself :-/
	if err != nil {
		// clientSpan.LogFields(log.String("event", "Tracer.Inject() failed"), log.Error(err))
		fmt.Println("[tzj] err: ", err)
	}
	return metadata.NewOutgoingContext(ctx, md)
}
