package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var meter metric.Meter

func HttpRequestMetrics() gin.HandlerFunc {
	meter = otel.Meter("processgo/api", metric.WithInstrumentationVersion("1.0.0"))

	reqHist, _ := meter.Float64Histogram(
		"http.server.duration", // 与 View 匹配；如不使用 View，可改为 "http_request_duration_seconds"
		metric.WithDescription("Histogram of response latency (seconds) of HTTP handlers."),
		metric.WithUnit("s"),
	)
	reqCounter, _ := meter.Int64Counter(
		"http.server.requests",
		metric.WithDescription("Counter of HTTP requests made."),
		metric.WithUnit("1"),
	)

	return func(c *gin.Context) {
		start := time.Now()

		// 统一使用路由模板，避免 label 爆炸
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		// Process request first
		c.Next()

		attrs := []attribute.KeyValue{
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.route", route),
			attribute.String("http.status_code", strconv.Itoa(c.Writer.Status())),
		}

		if reqCounter != nil {
			reqCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attrs...))
		}

		if reqHist != nil {
			reqHist.Record(
				c.Request.Context(),
				time.Since(start).Seconds(),
				metric.WithAttributes(attrs...),
			)
		}
	}
}
