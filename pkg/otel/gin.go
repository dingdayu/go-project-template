package otel

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func GinMiddleware() gin.HandlerFunc {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	return otelgin.Middleware(serviceName)
}

// Prometheus 监控访问点
func Prometheus(c *gin.Context) {
	// register promhttp.HandlerOpts DisableCompression
	promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
		DisableCompression: true,
		EnableOpenMetrics:  true,
	})).ServeHTTP(c.Writer, c.Request)
}
