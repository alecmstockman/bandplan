package middleware

import (
	requestlog "bandplan/src/logging"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipBucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type Limiter struct {
	mu          sync.Mutex
	buckets     map[string]*ipBucket
	nextCleanup time.Time
	interval    time.Duration
	burst       int
}

func NewRateLimiter(duration time.Duration, burstLimit int) *Limiter {

	limiter := Limiter{
		buckets:  make(map[string]*ipBucket),
		interval: duration,
		burst:    burstLimit,
	}
	return &limiter
}

func NewRegistrationLimiter() *Limiter {
	return NewRateLimiter(10*time.Second, 25)
}

func NewLoginLimiter() *Limiter {
	return NewRateLimiter(10*time.Second, 10)
}

func NewLoginEmailLimiter() *Limiter {
	return NewRateLimiter(5*time.Second, 30)
}

func (l *Limiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if !now.Before(l.nextCleanup) {
		for k, bucket := range l.buckets {
			if now.Sub(bucket.lastSeen) > 10*time.Minute {
				delete(l.buckets, k)
			}
		}
		l.nextCleanup = now.Add(time.Minute)
	}

	bucket, exists := l.buckets[key]
	if !exists {
		if len(l.buckets) >= 10_000 {
			return false
		}

		bucket = &ipBucket{
			limiter: rate.NewLimiter(rate.Every(l.interval), l.burst),
		}
		l.buckets[key] = bucket
	}
	bucket.lastSeen = now
	return bucket.limiter.Allow()
}

func (l *Limiter) Reset(key string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	delete(l.buckets, key)
}

func (l *Limiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := requestlog.GetClientIP(r.Context())
		if ip == "" {
			http.Error(w, "Client address unavailable", http.StatusBadRequest)
			return
		}

		if !l.Allow(ip) {
			w.Header().Set("Retry-After", "10")
			http.Error(w,
				"Too many requests. Please try again shortly.",
				http.StatusTooManyRequests,
			)
			return
		}

		next.ServeHTTP(w, r)
	})
}
