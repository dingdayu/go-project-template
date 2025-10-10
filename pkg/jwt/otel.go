package jwt

import "go.opentelemetry.io/otel"

var tracer = otel.Tracer("gitlab.keymedbio.com/platform/processgo/pkg/jwt")
