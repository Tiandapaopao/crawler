package ratelimit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	amount   int64
	max      int64
	lastTime int64
	rate     int64
	lock     sync.Mutex
}

func cur() int64 {
	return time.Now().Unix()
}

func New(rate int64, max int64) *RateLimiter {
	result := &RateLimiter{
		amount:   max,
		rate:     rate,
		max:      max,
		lastTime: cur(),
	}
	return result

}

func (r *RateLimiter) Pass() bool {
	r.lock.Lock()
	defer r.lock.Unlock()
	passed := cur() - r.lastTime
	amount := r.amount + passed*r.rate
	if amount > r.max {
		amount = r.max
	}

	if amount <= 0 {
		return false
	}
	amount--
	r.amount = amount
	r.lastTime = cur()
	return true
}
