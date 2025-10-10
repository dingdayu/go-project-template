package otel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime/debug"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/prometheus/client_golang/prometheus"
	prom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/resource"
)

// Options 控制初始化行为
type Options struct {
	Environment  string            // 可选：dev/staging/prod
	Endpoint     string            // 可选：例如 "otelcol.observability.svc:4318"（留空则走环境变量）
	Insecure     bool              // 集群内常用
	Headers      map[string]string // 需要鉴权时可传 {"Authorization":"Bearer xxx"}
	MetricPeriod time.Duration     // metrics push 周期，默认 10s
}

// Setup 初始化 OTEL：Trace + Metric + Log + Propagator
// 返回 shutdown 用于优雅退出
func Setup(ctx context.Context, opt Options) (shutdown func(context.Context) error, err error) {
	if opt.MetricPeriod <= 0 {
		opt.MetricPeriod = 10 * time.Second
	}

	// ---------- Propagator ----------
	// W3C 上下文传播（HTTP/RabbitMQ 都用它）
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// ---------- Resource（统一） ----------
	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),      // 允许通过 OTEL_SERVICE_NAME OTEL_RESOURCE_ATTRIBUTES 注入
		resource.WithTelemetrySDK(), // sdk 信息
		resource.WithAttributes(
			attribute.String("environment", opt.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("resource init: %w", err)
	}

	// ---------- Exporters ----------
	traceExp, err := newTraceExporter(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("trace exporter: %w", err)
	}
	metricExp, err := newMetricExporter(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("metric exporter: %w", err)
	}
	logExp, err := newLogExporter(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("log exporter: %w", err)
	}

	// ---------- Providers ----------
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(traceExp, sdktrace.WithBatchTimeout(1*time.Second)),
	)
	otel.SetTracerProvider(tp)

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(logExp)),
		// sdklog.WithResource(res), // 目前 LoggerProvider 不直接接收 Resource；resource 绑定在 exporter 端
	)
	global.SetLoggerProvider(lp)

	promExp, pErr := prom.New(prom.WithRegisterer(prometheus.DefaultRegisterer)) // pull exporter
	if pErr != nil {
		return nil, fmt.Errorf("prometheus exporter: %w", pErr)
	}

	// push exporter
	pushReader := sdkmetric.NewPeriodicReader(metricExp, sdkmetric.WithInterval(opt.MetricPeriod))
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(promExp),
		sdkmetric.WithReader(pushReader),
	)
	otel.SetMeterProvider(mp)

	// ---------- 汇总 shutdown ----------
	shutdown = func(c context.Context) error {
		var e error
		// 顺序建议：trace/metric/log 都可以，重要的是都要调用
		e = errors.Join(e, tp.Shutdown(c))
		e = errors.Join(e, mp.Shutdown(c))
		e = errors.Join(e, lp.Shutdown(c))
		return e
	}

	return shutdown, nil
}

// ---- exporters

func newTraceExporter(ctx context.Context, opt Options) (*otlptrace.Exporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithRetry(otlptracehttp.RetryConfig{ /* 默认重试 */ }),
	}
	if opt.Endpoint != "" {
		opts = append(opts, otlptracehttp.WithEndpoint(opt.Endpoint))
	}
	if opt.Insecure {
		opts = append(opts, otlptracehttp.WithInsecure())
	}
	if len(opt.Headers) > 0 {
		opts = append(opts, otlptracehttp.WithHeaders(opt.Headers))
	}
	return otlptracehttp.New(ctx, opts...)
}

func newMetricExporter(ctx context.Context, opt Options) (*otlpmetrichttp.Exporter, error) {
	opts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithRetry(otlpmetrichttp.RetryConfig{}),
	}
	if opt.Endpoint != "" {
		opts = append(opts, otlpmetrichttp.WithEndpoint(opt.Endpoint))
	}
	if opt.Insecure {
		opts = append(opts, otlpmetrichttp.WithInsecure())
	}
	if len(opt.Headers) > 0 {
		opts = append(opts, otlpmetrichttp.WithHeaders(opt.Headers))
	}
	return otlpmetrichttp.New(ctx, opts...)
}

func newLogExporter(ctx context.Context, opt Options) (*otlploghttp.Exporter, error) {
	opts := []otlploghttp.Option{
		otlploghttp.WithRetry(otlploghttp.RetryConfig{}),
	}
	if opt.Endpoint != "" {
		opts = append(opts, otlploghttp.WithEndpoint(opt.Endpoint))
	}
	if opt.Insecure {
		opts = append(opts, otlploghttp.WithInsecure())
	}
	if len(opt.Headers) > 0 {
		opts = append(opts, otlploghttp.WithHeaders(opt.Headers))
	}
	return otlploghttp.New(ctx, opts...)
}

func GetServiceName() string {
	if serviceName := os.Getenv("OTEL_SERVICE_NAME"); serviceName != "" {
		return serviceName
	}
	if serviceName := os.Getenv("APP_NAME"); serviceName != "" {
		return serviceName
	}
	if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
		return path.Base(bi.Main.Path)
	}
	return ""
}
