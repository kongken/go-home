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

// StringAttribute 创建字符串属性
func StringAttribute(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// IntAttribute 创建整数属性
func IntAttribute(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// Int64Attribute 创建 int64 属性
func Int64Attribute(key string, value int64) attribute.KeyValue {
	return attribute.Int64(key, value)
}

// BoolAttribute 创建布尔属性
func BoolAttribute(key string, value bool) attribute.KeyValue {
	return attribute.Bool(key, value)
}

// Float64Attribute 创建 float64 属性
func Float64Attribute(key string, value float64) attribute.KeyValue {
	return attribute.Float64(key, value)
}

// TraceUserRegister 追踪用户注册
func TraceUserRegister(ctx context.Context, userID, username string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "User.Register",
		StringAttribute("user.id", userID),
		StringAttribute("user.username", username),
	)
}

// TraceUserLogin 追踪用户登录
func TraceUserLogin(ctx context.Context, account string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "User.Login",
		StringAttribute("user.account", account),
	)
}

// TraceBlogCreate 追踪博客创建
func TraceBlogCreate(ctx context.Context, userID, blogID string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "Blog.Create",
		StringAttribute("user.id", userID),
		StringAttribute("blog.id", blogID),
	)
}

// TraceFeedCreate 追踪动态创建
func TraceFeedCreate(ctx context.Context, userID string, feedType int32) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, "Feed.Create",
		StringAttribute("user.id", userID),
		Int64Attribute("feed.type", int64(feedType)),
	)
}

// TraceCacheOperation 追踪缓存操作
func TraceCacheOperation(ctx context.Context, operation, cacheType, key string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, fmt.Sprintf("Cache.%s", operation),
		StringAttribute("cache.type", cacheType),
		StringAttribute("cache.key", key),
	)
}

// TraceDBOperation 追踪数据库操作
func TraceDBOperation(ctx context.Context, operation, table string) (context.Context, trace.Span) {
	return StartSpanWithAttributes(ctx, fmt.Sprintf("DB.%s", operation),
		StringAttribute("db.table", table),
	)
}