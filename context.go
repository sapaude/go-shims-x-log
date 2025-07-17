package log

import "context"

// contextKey 是一个私有类型，用于定义 Context 的键，避免冲突
type contextKey string

const (
    // RequestIDKey 用于在 Context 中存储请求 ID
    RequestIDKey contextKey = "request_id"
    // UserIDKey 用于在 Context 中存储用户 ID
    UserIDKey contextKey = "user_id"
    // TraceIDKey 用于在 Context 中存储链路追踪 ID
    TraceIDKey contextKey = "trace_id"
    // SpanIDKey 用于在 Context 中存储 Span ID
    SpanIDKey contextKey = "span_id"
    // CustomFieldsKey 用于在 Context 中存储一个 map[string]any，包含任意自定义字段
    CustomFieldsKey contextKey = "custom_fields"
)

// WithRequestID 将请求 ID 添加到 Context 中
func WithRequestID(ctx context.Context, reqID string) context.Context {
    return context.WithValue(ctx, RequestIDKey, reqID)
}

// WithUserID 将用户 ID 添加到 Context 中
func WithUserID(ctx context.Context, userID string) context.Context {
    return context.WithValue(ctx, UserIDKey, userID)
}

// WithTraceID 将 Trace ID 添加到 Context 中
func WithTraceID(ctx context.Context, traceID string) context.Context {
    return context.WithValue(ctx, TraceIDKey, traceID)
}

// WithSpanID 将 Span ID 添加到 Context 中
func WithSpanID(ctx context.Context, spanID string) context.Context {
    return context.WithValue(ctx, SpanIDKey, spanID)
}

type MetaData map[string]interface{}

// WithCustomField 将单个自定义字段添加到 Context 中。
// 如果 Context 中已有 CustomFieldsKey，则会更新或添加字段。
func WithCustomField(ctx context.Context, key string, value any) context.Context {
    fields, ok := ctx.Value(CustomFieldsKey).(MetaData)
    if !ok || fields == nil {
        fields = make(MetaData)
    } else {
        // 复制一份，避免修改原始 Context 中的 map
        newFields := make(MetaData, len(fields)+1)
        for k, v := range fields {
            newFields[k] = v
        }
        fields = newFields
    }
    fields[key] = value
    return context.WithValue(ctx, CustomFieldsKey, fields)
}

// GetRequestID 从 Context 中获取请求 ID
func GetRequestID(ctx context.Context) (string, bool) {
    val, ok := ctx.Value(RequestIDKey).(string)
    return val, ok
}

// GetUserID 从 Context 中获取用户 ID
func GetUserID(ctx context.Context) (string, bool) {
    val, ok := ctx.Value(UserIDKey).(string)
    return val, ok
}

// GetTraceID 从 Context 中获取 Trace ID
func GetTraceID(ctx context.Context) (string, bool) {
    val, ok := ctx.Value(TraceIDKey).(string)
    return val, ok
}

// GetSpanID 从 Context 中获取 Span ID
func GetSpanID(ctx context.Context) (string, bool) {
    val, ok := ctx.Value(SpanIDKey).(string)
    return val, ok
}

// GetCustomFields 从 Context 中获取所有自定义字段
func GetCustomFields(ctx context.Context) (MetaData, bool) {
    val, ok := ctx.Value(CustomFieldsKey).(MetaData)
    return val, ok
}
