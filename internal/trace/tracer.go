package trace

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Tracer 全局 tracer
	Tracer = otel.Tracer("go-home")
)

// StartSpan 开始 span
func StartSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name)
}

// StartSpanWithAttributes 开始 span 并添加属性
func StartSpanWithAttributes(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	ctx, span := Tracer.Start(ctx, name)
	span.SetAttributes(attrs...)
	return ctx, span
}

// RecordError 记录错误
func RecordError(span trace.Span, err error) {
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// AddEvent 添加事件
func AddEvent(span trace.Span, name string, attrs ...attribute.KeyValue) {
	span.AddEvent(name, trace.WithAttributes(attrs...))
}

// SetAttributes 设置属性
func SetAttributes(span trace.Span, attrs ...attribute.KeyValue) {
	span.SetAttributes(attrs...)
}

// Str 创建字符串属性
func Str(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// Int 创建整数属性
func Int(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// Int64 创建 int64 属性
func Int64(key string, value int64) attribute.KeyValue {
	return attribute.Int64(key, value)
}

// Bool 创建布尔属性
func Bool(key string, value bool) attribute.KeyValue {
	return attribute.Bool(key, value)
}

// Float64 创建 float64 属性
func Float64(key string, value float64) attribute.KeyValue {
	return attribute.Float64(key, value)
}

// UserRegister 追踪用户注册
func UserRegister(ctx context.Context, userID, username string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "User.Register",
		Str("user.id", userID),
		Str("user.username", username),
	)
}

// UserLogin 追踪用户登录
func UserLogin(ctx context.Context, account string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "User.Login",
		Str("user.account", account),
	)
}

// BlogCreate 追踪博客创建
func BlogCreate(ctx context.Context, userID, blogID string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "Blog.Create",
		Str("user.id", userID),
		Str("blog.id", blogID),
	)
}

// FeedCreate 追踪动态创建
func FeedCreate(ctx context.Context, userID string, feedType int32) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "Feed.Create",
		Str("user.id", userID),
		Int64("feed.type", int64(feedType)),
	)
}

// CacheOp 追踪缓存操作
func CacheOp(ctx context.Context, operation, cacheType, key string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, fmt.Sprintf("Cache.%s", operation),
		Str("cache.type", cacheType),
		Str("cache.key", key),
	)
}

// DBOp 追踪数据库操作
func DBOp(ctx context.Context, operation, table string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, fmt.Sprintf("DB.%s", operation),
		Str("db.table", table),
	)
}