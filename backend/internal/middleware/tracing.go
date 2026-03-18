package middleware

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HanzChrisrome/org-man-app/internal/config"
	"github.com/gin-gonic/gin"
)

type traceResponseWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *traceResponseWriter) Write(data []byte) (int, error) {
	w.body.Write(data)
	return w.ResponseWriter.Write(data)
}

func (w *traceResponseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

type payloadTraceRecord struct {
	TraceID        string            `json:"trace_id"`
	Timestamp      string            `json:"timestamp"`
	Method         string            `json:"method"`
	Path           string            `json:"path"`
	StatusCode     int               `json:"status_code"`
	DurationMS     int64             `json:"duration_ms"`
	RequestBody    string            `json:"request_body,omitempty"`
	ResponseBody   string            `json:"response_body,omitempty"`
	ResponseHeader map[string]string `json:"response_headers,omitempty"`
}

type payloadTraceErrorRecord struct {
	TraceID   string `json:"trace_id"`
	Timestamp string `json:"timestamp"`
	Method    string `json:"method,omitempty"`
	Path      string `json:"path,omitempty"`
	Stage     string `json:"stage"`
	Error     string `json:"error"`
}

func shouldTracePath(path string, excludes []string) bool {
	for _, pattern := range excludes {
		if pattern == "" {
			continue
		}

		if strings.HasSuffix(pattern, "*") {
			prefix := strings.TrimSuffix(pattern, "*")
			if strings.HasPrefix(path, prefix) {
				return false
			}
			continue
		}

		if path == pattern {
			return false
		}
	}

	return true
}

func clipPayload(data []byte, maxBytes int) string {
	if len(data) == 0 {
		return ""
	}

	if maxBytes <= 0 || len(data) <= maxBytes {
		return string(data)
	}

	return string(data[:maxBytes]) + "...[truncated]"
}

func selectHeaders(header http.Header, allowlist []string) map[string]string {
	if len(allowlist) == 0 {
		return nil
	}

	selected := make(map[string]string)
	for _, key := range allowlist {
		value := header.Get(key)
		if value != "" {
			selected[key] = value
		}
	}

	if len(selected) == 0 {
		return nil
	}

	return selected
}

func generateTraceID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return time.Now().UTC().Format("20060102T150405.000000000")
	}
	return hex.EncodeToString(b)
}

func logTraceError(traceLogger *log.Logger, traceID string, c *gin.Context, stage string, err error) {
	rec := payloadTraceErrorRecord{
		TraceID:   traceID,
		Timestamp: time.Now().UTC().Format(time.RFC3339Nano),
		Stage:     stage,
		Error:     err.Error(),
	}
	if c != nil && c.Request != nil {
		rec.Method = c.Request.Method
		rec.Path = c.Request.URL.Path
	}

	payload, marshalErr := json.Marshal(rec)
	if marshalErr != nil {
		traceLogger.Printf("{\"trace_id\":%q,\"stage\":%q,\"error\":%q}", traceID, stage, err.Error())
		return
	}

	traceLogger.Printf("%s", string(payload))
}

func PayloadTracer(cfg *config.Config) gin.HandlerFunc {
	outputs := []io.Writer{os.Stdout}
	if cfg.TraceFilePath != "" {
		dir := filepath.Dir(cfg.TraceFilePath)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0o755); err != nil {
				log.Printf("payload-trace mkdir failed for %s: %v", dir, err)
			}
		}

		file, err := os.OpenFile(cfg.TraceFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			log.Printf("payload-trace file open failed for %s: %v", cfg.TraceFilePath, err)
		} else {
			outputs = append(outputs, file)
		}
	}

	traceLogger := log.New(io.MultiWriter(outputs...), "payload-trace ", log.LstdFlags)

	return func(c *gin.Context) {
		if !cfg.EnablePayloadTrace || !shouldTracePath(c.Request.URL.Path, cfg.TraceExcludePaths) {
			c.Next()
			return
		}

		traceID := generateTraceID()
		c.Set("trace_id", traceID)
		c.Writer.Header().Set("X-Trace-ID", traceID)

		start := time.Now()

		var requestBody []byte
		if cfg.TraceRequestBody && c.Request != nil && c.Request.Body != nil {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				requestBody = bodyBytes
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			} else {
				logTraceError(traceLogger, traceID, c, "read_request_body", err)
				c.Request.Body = io.NopCloser(bytes.NewBuffer(nil))
			}
		}

		wrapped := &traceResponseWriter{ResponseWriter: c.Writer}
		c.Writer = wrapped

		c.Next()

		record := payloadTraceRecord{
			TraceID:    traceID,
			Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			StatusCode: c.Writer.Status(),
			DurationMS: time.Since(start).Milliseconds(),
		}

		if cfg.TraceRequestBody {
			record.RequestBody = clipPayload(requestBody, cfg.TraceMaxBodyBytes)
		}
		if cfg.TraceResponseBody {
			record.ResponseBody = clipPayload(wrapped.body.Bytes(), cfg.TraceMaxBodyBytes)
		}

		record.ResponseHeader = selectHeaders(c.Writer.Header(), cfg.TraceHeaders)

		payload, err := json.Marshal(record)
		if err != nil {
			logTraceError(traceLogger, traceID, c, "marshal_record", err)
			return
		}
		traceLogger.Printf("%s", string(payload))
	}
}
