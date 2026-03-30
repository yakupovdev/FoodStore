package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

type entry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
	banned   time.Time
}
type IPLimiter struct {
	limiters map[string]*entry
	mu       sync.Mutex
	r        rate.Limit
	b        int
}

func NewIPLimiter(ctx context.Context, r rate.Limit, b int) *IPLimiter {
	il := IPLimiter{
		limiters: make(map[string]*entry),
		r:        r,
		b:        b,
	}
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				il.cleanup(10 * time.Minute)
			}
		}
	}(ctx)
	return &il
}

func (i *IPLimiter) getLimiter(ip string) (*rate.Limiter, bool) {
	i.mu.Lock()
	defer i.mu.Unlock()

	e, ok := i.limiters[ip]
	if !ok {
		e = &entry{limiter: rate.NewLimiter(i.r, i.b)}
		i.limiters[ip] = e
	}
	e.lastSeen = time.Now()

	if time.Now().Before(e.banned) {
		return e.limiter, true
	}
	return e.limiter, false
}
func (i *IPLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter, banned := i.getLimiter(ip)

		if banned {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}
		if !limiter.Allow() {
			i.mu.Lock()
			i.limiters[ip].banned = time.Now().Add(5 * time.Minute)
			i.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})
			return
		}

		c.Next()
	}
}

func (i *IPLimiter) cleanup(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		i.mu.Lock()
		for ip, e := range i.limiters {
			if time.Since(e.lastSeen) > interval {
				delete(i.limiters, ip)
			}
		}
		i.mu.Unlock()
	}
}
