package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	RequestLogNamed    = "http_request"
	maxBodyLogSize     = 1024 // 设置请求体日志大小的最大阈值为 1024 bytes (1KB)
	maxResponseLogSize = 1024
)

// 一下响应文件不记录 response 到日志
var downloadableContentTypes = []string{
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", // Excel 文件
	"application/pdf",          // PDF 文件
	"application/octet-stream", // 通用二进制格式
	"image/jpeg",               // JPEG 图片
	"image/png",                // PNG 图片
	"image/gif",                // GIF 图片
	"image/bmp",                // BMP 图片
	"image/webp",               // WEBP 图片
	"text/csv",                 // CSV 文件
	"text/html",                // HTML 文件
	"text/javascript",          // JavaScript 文件
	"application/javascript",   // JavaScript 文件
	"text/css",                 // CSS 文件
	"font/ttf",                 // TrueType 字体
	"image/svg+xml",            // SVG 文件
	"application/zip",          // ZIP 文件
	"application/x-rar-compressed",
	"application/x-7z-compressed",
	"application/x-bzip2",
	"application/x-bzip",
	"application/x-gzip",
	// 添加更多文件类型
}

// WriterLog 处理跨域请求,支持options访问
func WriterLog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 准备相应日志
		bodyBuf := new(bytes.Buffer)
		_, _ = io.Copy(bodyBuf, c.Request.Body)
		body := bodyBuf.Bytes()
		c.Request.Body = io.NopCloser(bytes.NewReader(body))
		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		latency := time.Since(start)
		contentType := c.Writer.Header().Get("Content-Type")

		fs := []any{
			slog.Int("status", c.Writer.Status()),
			slog.String("ip", c.ClientIP()),
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.Int64("latency", latency.Milliseconds()),
			slog.String("user_agent", c.Request.UserAgent()),
		}
		// 当 query 不为空时记录
		if len(c.Request.URL.RawQuery) > 0 {
			fs = append(fs, slog.String("query", c.Request.URL.RawQuery))
		}
		// 检查响应类型黑名单
		if !isFileResponse(contentType) && blw.body.Len() <= maxResponseLogSize {
			fs = append(fs, slog.String("response", string(mask(blw.body.Bytes()))))
		}

		// Append error field if this is an erroneous request.
		if len(c.Errors) > 0 {
			fs = append(fs, slog.String("errors", c.Errors.String()))
		}
		// 非GET请求记录请求体，且只记录小于等于1KB的请求体
		if c.Request.Method != http.MethodGet && len(body) <= maxBodyLogSize {
			fs = append(fs, slog.Any("body", mask(body)))
		}
		// Writer X-Request-Id to log
		xRequestId := c.Request.Header.Get("X-Request-Id")
		if len(xRequestId) > 0 {
			fs = append(fs, slog.String("request_id", xRequestId))
		}

		logger.InfoContext(c.Request.Context(), c.Request.RequestURI, fs...)
	}
}

var maskDictionary = map[string]bool{"password": true}

func mask(body []byte) []byte {
	// 将 JSON 转换为 map
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return body
	}

	// 过滤敏感信息
	filterSensitiveData(data, maskDictionary)

	// 将过滤后的数据转换回 JSON 字符串
	filteredJSON, err := json.Marshal(data)
	if err != nil {
		return body
	}
	return filteredJSON
}

// 递归过滤敏感信息
func filterSensitiveData(data map[string]interface{}, maskDictionary map[string]bool) {
	for key, value := range data {
		// 检查是否为敏感字段
		if maskDictionary[key] {
			data[key] = "***"
		} else {
			// 递归处理嵌套的 map
			switch v := value.(type) {
			case map[string]interface{}:
				filterSensitiveData(v, maskDictionary)
			case []interface{}:
				for _, item := range v {
					if itemMap, ok := item.(map[string]interface{}); ok {
						filterSensitiveData(itemMap, maskDictionary)
					}
				}
			}
		}
	}
}

// 检查是否为文件响应
func isFileResponse(contentType string) bool {
	for _, fileType := range downloadableContentTypes {
		if strings.Contains(contentType, fileType) {
			return true
		}
	}
	return false
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
