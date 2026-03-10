package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestDuration HTTP 请求耗时
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "HTTP request latencies in seconds",
			Buckets: []float64{0.001, 0.01, 0.1, 0.5, 1, 2, 5},
		},
		[]string{"method", "path", "status"},
	)

	// HTTPRequestTotal HTTP 请求总数
	HTTPRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	// ActiveConnections 活跃连接数
	ActiveConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "active_connections",
			Help: "Number of active connections",
		},
	)

	// UserRegisteredTotal 注册用户总数
	UserRegisteredTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "user_registered_total",
			Help: "Total number of registered users",
		},
	)

	// BlogCreatedTotal 博客创建总数
	BlogCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "blog_created_total",
			Help: "Total number of blogs created",
		},
	)

	// FeedCreatedTotal 动态创建总数
	FeedCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "feed_created_total",
			Help: "Total number of feeds created",
		},
	)

	// CacheHitTotal 缓存命中总数
	CacheHitTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_hit_total",
			Help: "Total number of cache hits",
		},
		[]string{"cache_type"},
	)

	// CacheMissTotal 缓存未命中总数
	CacheMissTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cache_miss_total",
			Help: "Total number of cache misses",
		},
		[]string{"cache_type"},
	)
)

// CacheHit 记录缓存命中
func CacheHit(cacheType string) {
	CacheHitTotal.WithLabelValues(cacheType).Inc()
}

// CacheMiss 记录缓存未命中
func CacheMiss(cacheType string) {
	CacheMissTotal.WithLabelValues(cacheType).Inc()
}

// ActiveConnInc 增加活跃连接
func ActiveConnInc() {
	ActiveConnections.Inc()
}

// ActiveConnDec 减少活跃连接
func ActiveConnDec() {
	ActiveConnections.Dec()
}

// UserRegistered 增加注册用户
func UserRegistered() {
	UserRegisteredTotal.Inc()
}

// BlogCreated 增加博客创建
func BlogCreated() {
	BlogCreatedTotal.Inc()
}

// FeedCreated 增加动态创建
func FeedCreated() {
	FeedCreatedTotal.Inc()
}