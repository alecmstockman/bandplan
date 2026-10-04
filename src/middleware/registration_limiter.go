package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type ipBucket struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RegistrationLimiter struct {
	mu          sync.Mutex
	buckets     map[string]*ipBucket
	nextCleanup time.Time
}

func NewRegistrationLimiter() *RegistrationLimiter {
	return &RegistrationLimiter{
		buckets: make(map[string]*ipBucket),
	}
}

func (l *RegistrationLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()

	if !now.Before(l.nextCleanup) {
		for key, bucket := range l.buckets {
			if now.Sub(bucket.lastSeen) > 10*time.Minute {
				delete(l.buckets, key)
			}
		}
		l.nextCleanup = now.Add(time.Minute)
	}

	bucket, exists := l.buckets[ip]
	if !exists {
		if len(l.buckets) >= 10_000 {
			return false
		}

		bucket = &ipBucket{
			limiter: rate.NewLimiter(rate.Every(10*time.Second), 25),
		}
		l.buckets[ip] = bucket
	}
	bucket.lastSeen = now
	return bucket.limiter.Allow()
}

func (l *RegistrationLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Invalid client address", http.StatusBadRequest)
			return
		}

		if !l.allow(ip) {
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
